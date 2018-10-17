package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"

	authService "github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	channel struct {
		db  *factory.DB
		ctx context.Context

		usr authService.UserService
		evl EventService

		channel repository.ChannelRepository
		cmember repository.ChannelMemberRepository
		cview   repository.ChannelViewRepository
		message repository.MessageRepository

		sysmsgs types.MessageSet
	}

	ChannelService interface {
		With(ctx context.Context) ChannelService

		FindByID(channelID uint64) (*types.Channel, error)
		Find(filter *types.ChannelFilter) (types.ChannelSet, error)

		Create(channel *types.Channel) (*types.Channel, error)
		Update(channel *types.Channel) (*types.Channel, error)

		FindByMembership() (rval []*types.Channel, err error)
		FindMembers(channelID uint64) (types.ChannelMemberSet, error)

		InviteUser(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		DeleteMember(channelID uint64, memberIDs ...uint64) (err error)

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
		RecordView(channelID, userID, lastMessageID uint64) error
	}

	// channelSecurity interface {
	// 	CanRead(ch *types.Channel) bool
	// }
)

func Channel() ChannelService {
	return (&channel{
		usr: authService.DefaultUser,
		evl: DefaultEvent,
	}).With(context.Background())
}

func (svc *channel) With(ctx context.Context) ChannelService {
	db := repository.DB(ctx)
	return &channel{
		db:  db,
		ctx: ctx,

		usr: svc.usr.With(ctx),
		evl: svc.evl.With(ctx),

		channel: repository.Channel(ctx, db),
		cmember: repository.ChannelMember(ctx, db),
		cview:   repository.ChannelView(ctx, db),
		message: repository.Message(ctx, db),

		// System messages should be flushed at the end of each session
		sysmsgs: types.MessageSet{},
	}
}

func (svc *channel) FindByID(id uint64) (ch *types.Channel, err error) {
	ch, err = svc.channel.FindChannelByID(id)
	if err != nil {
		return
	}

	// if !svc.sec.ch.CanRead(ch) {
	// 	return nil, errors.New("Not allowed to access channel")
	// }

	return
}

func (svc *channel) Find(filter *types.ChannelFilter) (cc types.ChannelSet, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return cc, svc.db.Transaction(func() (err error) {
		if cc, err = svc.channel.FindChannels(filter); err != nil {
			return
		} else if err = svc.preloadExtras(cc); err != nil {
			return
		}

		return
	})
}

// preloadExtras pre-loads channel's members, views
func (svc *channel) preloadExtras(cc types.ChannelSet) (err error) {
	if err = svc.preloadMembers(cc); err != nil {
		return
	}

	if err = svc.preloadViews(cc); err != nil {
		return
	}

	return
}

func (svc *channel) preloadMembers(cc types.ChannelSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if mm, err := svc.cmember.Find(&types.ChannelMemberFilter{ComembersOf: userID}); err != nil {
		return err
	} else {
		cc.Walk(func(ch *types.Channel) error {
			ch.Members = mm.MembersOf(ch.ID)
			return nil
		})
	}

	return nil
}

func (svc *channel) preloadViews(cc types.ChannelSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if vv, err := svc.cview.Find(&types.ChannelViewFilter{UserID: userID}); err != nil {
		return err
	} else {
		cc.Walk(func(ch *types.Channel) error {
			ch.View = vv.FindByChannelId(ch.ID)
			return nil
		})
	}

	return nil
}

// FindMembers loads all members (and full users) for a specific channel
func (svc *channel) FindMembers(channelID uint64) (out types.ChannelMemberSet, err error) {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	// @todo [SECURITY] check if we can return members on this channel
	_ = channelID
	_ = userID

	return out, svc.db.Transaction(func() (err error) {
		out, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID})
		if err != nil {
			return err
		}

		if uu, err := svc.usr.Find(nil); err != nil {
			return err
		} else {
			return out.Walk(func(member *types.ChannelMember) error {
				member.User = uu.FindById(member.UserID)
				return nil
			})
		}
	})
}

// Returns all channels with membership info
func (svc *channel) FindByMembership() (cc []*types.Channel, err error) {
	return cc, svc.db.Transaction(func() error {
		var chMemberId = repository.Identity(svc.ctx)

		var mm []*types.ChannelMember

		if mm, err = svc.cmember.Find(&types.ChannelMemberFilter{MemberID: chMemberId}); err != nil {
			return err
		}

		if cc, err = svc.channel.FindChannels(nil); err != nil {
			return err
		} else if err = svc.preloadExtras(cc); err != nil {
			return err
		}

		for _, m := range mm {
			for _, c := range cc {
				if c.ID == m.ChannelID {
					c.Member = m
				}
			}
		}

		return nil
	})
}

func (svc *channel) Create(in *types.Channel) (out *types.Channel, err error) {
	// @todo: [SECURITY] permission check if user can add channel

	return out, svc.db.Transaction(func() (err error) {
		var msg *types.Message

		// @todo get organisation from somewhere
		var organisationID uint64 = 0

		var chCreatorID = repository.Identity(svc.ctx)

		// @todo [SECURITY] check if channel topic can be set
		if in.Topic != "" && false {
			return errors.New("Not allowed to set channel topic")
		}

		// @todo [SECURITY] check if user can create public channels
		if in.Type == types.ChannelTypePublic && false {
			return errors.New("Not allowed to create public channels")
		}

		// @todo [SECURITY] check if user can create private channels
		if in.Type == types.ChannelTypePrivate && false {
			return errors.New("Not allowed to create public channels")
		}

		// @todo [SECURITY] check if user can create private channels
		if in.Type == types.ChannelTypeGroup && false {
			return errors.New("Not allowed to create group channels")
		}

		// This is a fresh channel, just copy values
		out = &types.Channel{
			Name:           in.Name,
			Topic:          in.Topic,
			Type:           in.Type,
			OrganisationID: organisationID,
			CreatorID:      chCreatorID,
		}

		// Save the channel
		if out, err = svc.channel.CreateChannel(out); err != nil {
			return
		}

		// Join current user as an member & owner
		_, err = svc.cmember.Create(&types.ChannelMember{
			ChannelID: out.ID,
			UserID:    chCreatorID,
			Type:      types.ChannelMembershipTypeOwner,
		})

		if err != nil {
			// Could not add member
			return
		}

		// When broadcasting, make sure we let the subscribers know
		// that creator is also a member...
		out.Members = append(out.Members, chCreatorID)

		svc.evl.Join(chCreatorID, out.ID)

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		svc.scheduleSystemMessage(
			out,
			"@%d created new %s channel, topic is: %s",
			chCreatorID,
			"<PRIVATE-OR-PUBLIC>",
			"<TOPIC>")

		_ = msg
		if err != nil {
			// Message creation failed
			return
		}

		svc.flushSystemMessages()

		return svc.sendChannelEvent(out)
	})
}

func (svc *channel) Update(ch *types.Channel) (out *types.Channel, err error) {
	return out, svc.db.Transaction(func() (err error) {
		// @todo [SECURITY] can user access this channel?
		if out, err = svc.channel.FindChannelByID(ch.ID); err != nil {
			return
		}

		if out.ArchivedAt != nil {
			return errors.New("Not allowed to edit archived channels")
		} else if out.DeletedAt != nil {
			return errors.New("Not allowed to edit deleted channels")
		}

		if out.Type != ch.Type {
			// @todo [SECURITY] check if user can create public channels
			if ch.Type == types.ChannelTypePublic && false {
				return errors.New("Not allowed to change type of this channel to public")
			}

			// @todo [SECURITY] check if user can create private channels
			if ch.Type == types.ChannelTypePrivate && false {
				return errors.New("Not allowed to change type of this channel to private")
			}

			// @todo [SECURITY] check if user can create group channels
			if ch.Type == types.ChannelTypeGroup && false {
				return errors.New("Not allowed to change type of this channel to group")
			}
		}

		var chUpdatorId = repository.Identity(svc.ctx)

		// Copy values
		if out.Name != ch.Name {
			// @todo [SECURITY] can we change channel's name?
			if false {
				return errors.New("Not allowed to rename channel")
			} else {
				svc.scheduleSystemMessage(ch, "@%d renamed channel %s (was: %s)", chUpdatorId, out.Name, ch.Name)
			}
			out.Name = ch.Name
		}

		if out.Topic != ch.Topic && true {
			// @todo [SECURITY] can we change channel's topic?
			if false {
				return errors.New("Not allowed to change channel topic")
			} else {
				svc.scheduleSystemMessage(ch, "@%d changed channel topic: %s (was: %s)", chUpdatorId, out.Topic, ch.Topic)
			}

			out.Topic = ch.Topic
		}

		// Save the updated channel
		if out, err = svc.channel.UpdateChannel(ch); err != nil {
			return
		}

		svc.flushSystemMessages()

		return svc.sendChannelEvent(out)
	})
}

func (svc *channel) Delete(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user delete this channel?

		if ch.DeletedAt != nil {
			return errors.New("Channel already deleted")
		} else {
			now := time.Now()
			ch.DeletedAt = &now
		}

		svc.scheduleSystemMessage(ch, "@%d deleted this channel", userID)

		if err = svc.channel.DeleteChannelByID(id); err != nil {
			return
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Recover(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user recover this channel?

		if ch.DeletedAt == nil {
			return errors.New("Channel not deleted")
		}

		svc.scheduleSystemMessage(ch, "@%d recovered this channel", userID)

		if err = svc.channel.UnarchiveChannelByID(id); err != nil {
			return
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Archive(id uint64) error {
	return svc.db.Transaction(func() (err error) {

		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user archive this channel?

		if ch.ArchivedAt != nil {
			return errors.New("Channel already archived")
		}

		svc.scheduleSystemMessage(ch, "@%d archived this channel", userID)

		if err = svc.channel.ArchiveChannelByID(id); err != nil {
			return
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Unarchive(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user unarchive this channel?

		if ch.ArchivedAt == nil {
			return errors.New("Channel not archived")
		}

		svc.scheduleSystemMessage(ch, "@%d unarchived this channel", userID)

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) InviteUser(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	out = types.ChannelMemberSet{}

	// @todo [SECURITY] can user access this channel?
	if ch, err = svc.channel.FindChannelByID(channelID); err != nil {
		return
	}

	// @todo [SECURITY] can user add members to this channel?

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		users, err := svc.usr.Find(nil)
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			user := users.FindById(memberID)
			if user == nil {
				return errors.New("Unexisting user")
			}

			if e := existing.FindByUserId(memberID); e != nil {
				// Already a member/invited
				e.User = user
				out = append(out, e)
				continue
			}

			svc.scheduleSystemMessage(ch, "@%d invited @%d to the channel", userID, memberID)

			member := &types.ChannelMember{
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeInvitee,
			}

			if member, err = svc.cmember.Create(member); err != nil {
				return err
			}

			out = append(out, member)
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	out = types.ChannelMemberSet{}

	// @todo [SECURITY] can user access this channel?
	if ch, err = svc.channel.FindChannelByID(channelID); err != nil {
		return
	}

	// @todo [SECURITY] can user add members to this channel?

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		users, err := svc.usr.Find(nil)
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			var exists bool

			user := users.FindById(memberID)
			if user == nil {
				return errors.New("Unexisting user")
			}

			if e := existing.FindByUserId(memberID); e != nil {
				if e.Type != types.ChannelMembershipTypeInvitee {
					e.User = user
					out = append(out, e)
					continue
				} else {
					exists = true
				}
			}

			if !exists {
				if userID == memberID {
					svc.scheduleSystemMessage(ch, "@%d joined", memberID)
				} else {
					svc.scheduleSystemMessage(ch, "@%d added @%d to the channel", userID, memberID)
				}
			}

			member := &types.ChannelMember{
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeOwner,
				User:      user,
			}

			if exists {
				member, err = svc.cmember.Update(member)
			} else {
				member, err = svc.cmember.Create(member)
			}

			svc.evl.Join(memberID, channelID)

			if err != nil {
				return err
			}

			out = append(out, member)
		}

		// Push channel to all members
		if err = svc.sendChannelEvent(ch); err != nil {
			return
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) DeleteMember(channelID uint64, memberIDs ...uint64) (err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	// @todo [SECURITY] can user access this channel?
	if ch, err = svc.channel.FindChannelByID(channelID); err != nil {
		return
	}

	// @todo [SECURITY] can user remove members from this channel?

	return svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		for _, memberID := range memberIDs {
			if existing.FindByUserId(memberID) == nil {
				// Not really a member...
				continue
			}

			if userID == memberID {
				svc.scheduleSystemMessage(ch, "@%d parted", memberID)
			} else {
				svc.scheduleSystemMessage(ch, "@%d kicked @%d out", userID, memberID)
			}

			if err = svc.cmember.Delete(channelID, memberID); err != nil {
				return err
			}

			svc.evl.Part(memberID, channelID)
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) RecordView(channelID, userID, lastMessageID uint64) error {
	return svc.db.Transaction(func() (err error) {
		return svc.cview.Record(channelID, userID, lastMessageID, 0)
	})
}

func (svc *channel) scheduleSystemMessage(ch *types.Channel, format string, a ...interface{}) {
	svc.sysmsgs = append(svc.sysmsgs, &types.Message{
		ChannelID: ch.ID,
		Message:   fmt.Sprintf(format, a...),
		Type:      types.MessageTypeChannelEvent,
	})
}

// Flushes sys message stack, stores them into repo & pushes them into event loop
func (svc *channel) flushSystemMessages() (err error) {
	defer func() {
		svc.sysmsgs = types.MessageSet{}
	}()

	return svc.sysmsgs.Walk(func(msg *types.Message) error {
		if msg, err = svc.message.CreateMessage(msg); err != nil {
			return err
		} else {
			return svc.evl.Message(msg)
		}
	})
}

// Sends channel event
func (svc *channel) sendChannelEvent(ch *types.Channel) (err error) {
	if ch.DeletedAt == nil && ch.ArchivedAt == nil {
		// Looks like a valid channel

		// Preload members, if needed
		if len(ch.Members) == 0 {
			if mm, err := svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: ch.ID}); err != nil {
				return err
			} else {
				ch.Members = mm.MembersOf(ch.ID)
			}
		}
	}

	if err = svc.evl.Channel(ch); err != nil {
		return
	}

	return nil
}

//// @todo temp location, move this somewhere else
//type (
//	nativeChannelSec struct {
//		rpo struct {
//			ch nativeChannelSecChRepo
//		}
//	}
//
//	nativeChannelSecChRepo interface {
//		FindMember(channelId uint64, userId uint64) (*types.User, error)
//	}
//)
//
//func ChannelSecurity(chRpo nativeChannelSecChRepo) channelSecurity {
//	var sec = &nativeChannelSec{}
//
//	sec.rpo.ch = chRpo
//
//	return sec
//}
//
//// Current user can read the channel if he is a member
//func (sec nativeChannelSec) CanRead(ch *types.Channel) bool {
//	// @todo check if channel is public?
//
//	var currentUserID = repository.Identity(svc.ctx)
//
//	user, err := sec.rpo.FindMember(ch.ID, currentUserID)
//
//	return err != nil && user.Valid()
//}

var _ ChannelService = &channel{}
