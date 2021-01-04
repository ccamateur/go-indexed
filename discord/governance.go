package discord

import (
	"fmt"
	"log"
	"math/big"

	"github.com/bonedaddy/dgc"
	"github.com/bwmarrin/discordgo"
)

func (c *Client) governanceCurrentProposalsHandler(ctx *dgc.Ctx) {
	if !ctx.Command.RateLimiter.NotifyExecution(ctx) {
		return
	}
	gov, err := c.bc.GovernorAlpha()
	if err != nil {
		log.Println("governor alpha fetch failure: ", err)
		ctx.RespondText("failed to get governance contract handler")
		return
	}
	count, err := gov.ProposalCount(nil)
	if err != nil {
		log.Println("failed to get proposal count: ", err)
		ctx.RespondText("failed to get proposal count")
		return
	}
	ctx.RespondText(fmt.Sprint("number of proposals submitted: ", count))
}

func (c *Client) governanceProposalInfoHandler(ctx *dgc.Ctx) {
	arguments := ctx.Arguments
	number, err := arguments.Get(0).AsInt64()
	if err != nil {
		ctx.RespondText("invalid argument given must be a number")
		return
	}
	if number == 0 {
		ctx.RespondText("number must be 1 or higher")
		return
	}
	if !ctx.Command.RateLimiter.NotifyExecution(ctx) {
		return
	}
	gov, err := c.bc.GovernorAlpha()
	if err != nil {
		log.Println("governor alpha fetch failure: ", err)
		ctx.RespondText("failed to get governance contract handler")
		return
	}
	proposal, err := gov.Proposals(nil, big.NewInt(number))
	if err != nil {
		log.Println("failed to get proposal: ", err)
		ctx.RespondText("failed to get proposal")
		return
	}
	/*
			Id           *big.Int
		Proposer     common.Address
		Eta          *big.Int
		StartBlock   *big.Int
		EndBlock     *big.Int
		ForVotes     *big.Int
		AgainstVotes *big.Int
		Canceled     bool
		Executed     bool
	*/
	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Proposal Overview",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "ID",
				Value: proposal.Id.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "Proposer",
				Value: proposal.Proposer.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "ETA",
				Value: proposal.Eta.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "StartBlock",
				Value: proposal.StartBlock.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "EndBlock",
				Value: proposal.EndBlock.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "ForVotes",
				Value: proposal.ForVotes.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "AgainstVotes",
				Value: proposal.AgainstVotes.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "Canceled",
				Value: fmt.Sprint(proposal.Canceled),
			},
			&discordgo.MessageEmbedField{
				Name:  "Executed",
				Value: fmt.Sprint(proposal.Executed),
			},
		},
		Color: 0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://pbs.twimg.com/profile_images/1342395531318976518/kIv5abLc_400x400.jpg",
		},
	})
}