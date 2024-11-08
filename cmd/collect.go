package cmd

import (
	"fmt"
	"log"
	"strconv"

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
		for i, miner := range miners {
			fmt.Printf("Miner: %d/%d: %s\n", i+1, len(miners), miner.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
}
