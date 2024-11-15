package cmd

import (
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/sdk"
	"github.com/glifio/go-pools/terminate"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TestCase struct {
	SectorSizeGiB uint64
	Days          uint64
	RatioVerified float64
}

// penaltyOnPledgeCmd represents the penalty-on-pledge command
var penaltyOnPledgeCmd = &cobra.Command{
	Use:   "penalty-on-pledge <epoch>",
	Short: "Calculate the termination penalties for 1000 FIL in various scenarios",
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
		filToPledge := new(big.Int).Mul(big.NewInt(1000), constants.WAD)

		testcases := []TestCase{
			{SectorSizeGiB: 32, Days: 180, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 180, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 180, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 180, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 360, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 360, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 360, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 360, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 540, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 540, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 540, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 540, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 720, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 720, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 720, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 720, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 900, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 900, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 900, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 900, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 1080, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 1080, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 1080, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 1080, RatioVerified: 1.0},
			{SectorSizeGiB: 32, Days: 1260, RatioVerified: 0.0},
			{SectorSizeGiB: 32, Days: 1260, RatioVerified: 1.0},
			{SectorSizeGiB: 64, Days: 1260, RatioVerified: 0.0},
			{SectorSizeGiB: 64, Days: 1260, RatioVerified: 1.0},
		}
		for _, testcase := range testcases {
			sectorSize := testcase.SectorSizeGiB * 1073741824
			activation := ts.Height()
			expiration := activation + abi.ChainEpoch(testcase.Days*2880)

			cost, penalty, sectors, pledge, err := terminate.TermPenaltyOnPledge(ctx, lotusClient, ts,
				filToPledge, sectorSize, activation, expiration, testcase.RatioVerified)
			if err != nil {
				log.Fatal(err)
			}
			costFIL, _ := util.ToFIL(cost).Float64()
			penaltyFIL, _ := util.ToFIL(penalty).Float64()
			pct := penaltyFIL / costFIL * 100
			fmt.Printf("%dGiB, %d days, %0.1f%% verified: %0.2f%% penalty"+
				" (%0.1f FIL / %0.1f FIL, %d sectors, pledge: %0.2f FIL/sector)\n",
				testcase.SectorSizeGiB,
				testcase.Days,
				testcase.RatioVerified*100,
				pct,
				penaltyFIL,
				costFIL,
				sectors,
				util.ToFIL(pledge),
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(penaltyOnPledgeCmd)
}
