package cmd

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/glifio/go-pools/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Extract list of miners and sector data for a specific epoch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		epoch, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		nodeDialAddr := viper.GetString("lotus_addr")
		nodeToken := viper.GetString("lotus_token")

		lotusClient, closer, err := sdk.ConnectLotusClient(nodeDialAddr, nodeToken)
		if err != nil {
			log.Fatalf("failed to connect to lotus client: %v", err)
		}
		defer closer()

		ts, err := lotusClient.ChainGetTipSetByHeight(ctx, abi.ChainEpoch(epoch), types.EmptyTSK)
		if err != nil {
			log.Fatalf("could not get tipset: %v\n", err)
		}

		miners, err := lotusClient.StateListMiners(ctx, ts.Key())
		if err != nil {
			log.Fatalf("error listing miners: %v\n", err)
		}
		slices.SortFunc(miners, func(a, b address.Address) int {
			aMinerID, _ := strconv.ParseUint(a.String()[1:], 10, 64)
			bMinerID, _ := strconv.ParseUint(b.String()[1:], 10, 64)
			return cmp.Compare(aMinerID, bMinerID)
		})
		count := 0
		for i, miner := range miners {
			sectorCount, err := lotusClient.StateMinerSectorCount(ctx, miner, ts.Key())
			if err != nil {
				log.Fatalf("error getting sector count for %v: %v", miner, err)
			}
			if sectorCount.Active > 0 {
				count++
				fmt.Printf("#%d: %d/%d: %s (%d sectors)\n", count, i+1, len(miners),
					miner.String(), sectorCount.Active)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
}
