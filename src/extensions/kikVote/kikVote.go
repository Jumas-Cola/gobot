package kikVote

import (
	"context"
	"fmt"
	"gobot/src/db"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"
)

type KikVote struct {
	collection *mongo.Collection
	maxVotes   int
}

type VoteStatusType int32

const (
	Pending VoteStatusType = 0
	Kik     VoteStatusType = 1
	Forgive VoteStatusType = 2
)

type ChatVote struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ChatId       int64              `bson:"chatId"`
	MsgId        int                `bson:"msgId"`
	KikUserId    int64              `bson:"kikUserId"`
	Status       VoteStatusType     `bson:"status"`
	VotesFor     []int64            `bson:"votesFor"`
	VotesAgainst []int64            `bson:"votesAgainst"`
}

func GetExtension() KikVote {
	return KikVote{}
}

func (k KikVote) RegisterHandlers(b *telebot.Bot) []telebot.Command {
	cmds := []telebot.Command{
		{Text: "kik", Description: "Vote for kik"},
	}

	k.maxVotes = 5

	k.collection = db.MongoDbClient.Database("kikVote").Collection("ChatVotes")

	b.Handle("/kik", k.handleKik)

	b.Handle(&telebot.InlineButton{Unique: "kikVoteForBtn"}, k.handleKikVoteBtn)

	b.Handle(&telebot.InlineButton{Unique: "kikVoteAgainstBtn"}, k.handleKikVoteBtn)

	return cmds
}

func (k KikVote) handleKik(c telebot.Context) error {
	message := c.Message()
	isReply := message.IsReply()
	if !isReply {
		return c.Send("@" + message.Sender.Username + ", команда должна быть ответом на другое сообщение.")
	}

	newItem := ChatVote{
		KikUserId:    message.ReplyTo.Sender.ID,
		ChatId:       c.Chat().ID,
		VotesFor:     []int64{message.Sender.ID},
		VotesAgainst: []int64{},
	}
	res, err := k.collection.InsertOne(context.TODO(), newItem)
	if err != nil {
		log.Fatal(err)
	}

	inlineKeys := [][]telebot.InlineButton{k.makeVoteBtns(newItem)}

	inlineKeyboard := &telebot.ReplyMarkup{}
	inlineKeyboard.InlineKeyboard = inlineKeys

	bot := c.Bot()
	originalMessage := message.ReplyTo
	sentMsg, sentErr := bot.Send(c.Chat(), "@"+message.Sender.Username+
		" хочет кикнуть @"+
		originalMessage.Sender.Username+
		" из чата. Согласны?", inlineKeyboard)
	if sentErr != nil {
		log.Fatal(sentErr)
	}

	update := bson.M{
		"$set": bson.M{
			"msgId": sentMsg.ID,
		},
	}

	_, err = k.collection.UpdateByID(context.TODO(), res.InsertedID, update)
	if err != nil {
		log.Fatal(err)
	}

	return sentErr
}

func (k KikVote) handleKikVoteBtn(c telebot.Context) error {
	// TODO: проверять статус и роли пользователей

	message := c.Message()
	chat := c.Chat()
	callback := c.Callback()

	vote := k.getVote(*message, *chat)

	slog.Warn(fmt.Sprintf("%v", vote.Status))
	slog.Warn(fmt.Sprintf("%v", Pending))

	if vote.Status != Pending {
		return c.Respond(&telebot.CallbackResponse{})
	}

	update := bson.M{}

	for k, v := range vote.VotesAgainst {
		if v == callback.Sender.ID {
			vote.VotesAgainst = append(vote.VotesAgainst[:k], vote.VotesAgainst[k+1:]...)
			break
		}
	}

	for k, v := range vote.VotesFor {
		if v == callback.Sender.ID {
			vote.VotesFor = append(vote.VotesFor[:k], vote.VotesFor[k+1:]...)
			break
		}
	}

	switch callback.Unique {
	case "kikVoteForBtn":
		vote.VotesFor = append(vote.VotesFor, callback.Sender.ID)
	case "kikVoteAgainstBtn":
		vote.VotesAgainst = append(vote.VotesAgainst, callback.Sender.ID)
	}

	update = bson.M{
		"$set": bson.M{
			"votesFor":     vote.VotesFor,
			"votesAgainst": vote.VotesAgainst,
		},
	}

	_, err := k.collection.UpdateByID(context.TODO(), vote.ID, update)
	if err != nil {
		log.Fatal(err)
	}

	status := k.checkKikVoteFinish(vote)
	switch status {
	case Kik:
		k.kik(c, vote, message)
	case Forgive:
		k.forgive(c, vote, message)
	}

	inlineKeys := [][]telebot.InlineButton{k.makeVoteBtns(vote)}

	newMarkup := &telebot.ReplyMarkup{InlineKeyboard: inlineKeys}

	c.Edit(c.Message().Text, newMarkup)

	return c.Respond(&telebot.CallbackResponse{})
}

func (k KikVote) makeVoteBtns(vote ChatVote) []telebot.InlineButton {
	kikVoteForBtn := telebot.InlineButton{
		Unique: "kikVoteForBtn",
		Text:   fmt.Sprintf("☠️ Кикнуть (%d/%d)", len(vote.VotesFor), k.maxVotes),
	}
	kikVoteAgainstBtn := telebot.InlineButton{
		Unique: "kikVoteAgainstBtn",
		Text:   fmt.Sprintf("👼 Простить (%d/%d)", len(vote.VotesAgainst), k.maxVotes),
	}

	return []telebot.InlineButton{kikVoteForBtn, kikVoteAgainstBtn}
}

func (k KikVote) checkKikVoteFinish(vote ChatVote) VoteStatusType {
	if len(vote.VotesFor) >= k.maxVotes {
		return Kik
	} else if len(vote.VotesAgainst) >= k.maxVotes {
		return Forgive
	} else {
		return Pending
	}
}

func (k KikVote) getVote(message telebot.Message, chat telebot.Chat) ChatVote {
	filter := bson.M{"msgId": message.ID, "chatId": chat.ID}

	var res ChatVote

	err := k.collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (k KikVote) kik(c telebot.Context, vote ChatVote, message *telebot.Message) error {
	// TODO: добавить логику исключения
	update := bson.M{
		"$set": bson.M{
			"status": Kik,
		},
	}

	_, err := k.collection.UpdateByID(context.TODO(), vote.ID, update)
	if err != nil {
		log.Fatal(err)
	}

	return c.Send("@" + message.Sender.Username + " был кикнут из чата.")
}

func (k KikVote) forgive(c telebot.Context, vote ChatVote, message *telebot.Message) error {
	update := bson.M{
		"$set": bson.M{
			"status": Forgive,
		},
	}

	_, err := k.collection.UpdateByID(context.TODO(), vote.ID, update)
	if err != nil {
		log.Fatal(err)
	}

	return c.Send("@" + message.Sender.Username + " был прощён.")
}
