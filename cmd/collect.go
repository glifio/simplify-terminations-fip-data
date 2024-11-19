package cmd

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/glifio/go-pools/econ"
	"github.com/glifio/go-pools/sdk"
	"github.com/glifio/go-pools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect <epoch> [--save-sectors-dir <dir>]",
	Short: "Extract list of miners and sector data for a specific epoch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		epoch, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		progress, err := cmd.Flags().GetBool("progress")
		if err != nil {
			log.Fatal(err)
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatal(err)
		}

		outputCSV, err := cmd.Flags().GetBool("csv")
		if err != nil {
			log.Fatal(err)
		}

		saveSectorsSubdir := ""
		saveSectorsDir, err := cmd.Flags().GetString("save-sectors-dir")
		if err != nil {
			log.Fatal(err)
		}
		if saveSectorsDir != "" {
			saveSectorsSubdir = fmt.Sprintf("%s/%d", saveSectorsDir, epoch)
			os.MkdirAll(saveSectorsSubdir, 0750)
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

		start := time.Now()
		if progress {
			log.Println("Getting list of miners...")
		}
		miners, err := lotusClient.StateListMiners(ctx, ts.Key())
		if err != nil {
			log.Fatalf("error listing miners: %v\n", err)
		}
		if progress {
			elapsed := time.Since(start).Seconds()
			log.Printf("Found %d miners in %0.1f seconds\n", len(miners), elapsed)
		}
		slices.SortFunc(miners, func(a, b address.Address) int {
			aMinerID, _ := strconv.ParseUint(a.String()[1:], 10, 64)
			bMinerID, _ := strconv.ParseUint(b.String()[1:], 10, 64)
			return cmp.Compare(aMinerID, bMinerID)
		})
		w := csv.NewWriter(os.Stdout)
		if err := w.Write([]string{
			"Miner",
			"Epoch",
			"TotalBalance",
			"TotalBalanceFIL",
			"AvailableBalance",
			"AvailableBalanceFIL",
			"VestingFunds",
			"VestingFundsFIL",
			"InitialPledge",
			"InitialPledgeFIL",
			"FeeDebt",
			"FeeDebtFIL",
			"TerminationFee",
			"TerminationFeeFIL",
			"AvgTerminationFeePerPledge",
			"AvgTerminationFeePerPledgeFIL",
			"LiveSectors",
			"FaultySectors",
		}); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}

		count := 0
		for i, miner := range miners {
			start := time.Now()
			sectorCount, err := lotusClient.StateMinerSectorCount(ctx, miner, ts.Key())
			if err != nil {
				log.Fatalf("error getting sector count for %v: %v", miner, err)
			}

			if sectorCount.Live > 0 {
				sectors, err := econ.AllSectors(ctx, lotusClient, miner, ts)
				if err != nil {
					log.Fatalf("error getting sectors for %v: %v", miner, err)
				}
				sample := bitfield.NewFromSet(sectors)

				if saveSectorsSubdir != "" {
					filename := fmt.Sprintf("%s/%s.csv", saveSectorsSubdir, miner)
					sectorsFile, err := os.Create(filename)

					if err != nil {
						log.Fatalf("error creating file: %v", filename)
					}
					defer sectorsFile.Close()
					wSectors := csv.NewWriter(sectorsFile)
					if err := wSectors.Write([]string{
						"Miner",
						"Epoch",
						"SectorNumber",
						"SealProof",
						"SealedCID",
						"DealIDs",
						"Activation",
						"Expiration",
						"DealWeight",
						"VerifiedDealWeight",
						"InitialPledge",
						"ExpectedDayReward",
						"ExpectedStoragePledge",
						"PowerBaseEpoch",
						"ReplacedDayReward",
						"SectorKeyCID",
						"Flags",
					}); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}

					sectorInfos, err := lotusClient.StateMinerSectors(ctx, miner, &sample, ts.Key())
					if err != nil {
						log.Fatalf("error getting sectors for %v: %v", miner, err)
					}
					for _, sector := range sectorInfos {
						fmt.Printf("Jim sector: %+v\n", sector)
						if err := wSectors.Write([]string{
							miner.String(),           // Miner
							fmt.Sprintf("%d", epoch), // Epoch
						}); err != nil {
							log.Fatalln("error writing record to csv:", err)
						}
					}

					wSectors.Flush()
					if err := wSectors.Error(); err != nil {
						log.Fatal(err)
					}
				}

				res, err := econ.TerminateSectors(ctx, lotusClient, miner, &sample, ts)
				if err != nil {
					log.Fatalf("error terminating sectors for %v: %v", miner, err)
				}
				elapsed := time.Since(start).Seconds()
				count++
				if progress {
					log.Printf("#%d: %d/%d: %s (%d/%d active/live sectors, %0.1fs)\n", count, i+1, len(miners),
						miner.String(), sectorCount.Active, sectorCount.Live, elapsed)
				}
				if debug {
					log.Printf("%+v\n", res)
				}
				if outputCSV {
					if err := w.Write([]string{
						miner.String(),                                                   // Miner
						fmt.Sprintf("%d", epoch),                                         // Epoch
						res.TotalBalance.String(),                                        // TotalBalance
						fmt.Sprintf("%0.3f", util.ToFIL(res.TotalBalance)),               // TotalBalanceFIL
						res.AvailableBalance.String(),                                    // AvailableBalance
						fmt.Sprintf("%0.3f", util.ToFIL(res.AvailableBalance)),           // AvailableBalanceFIL
						res.VestingFunds.String(),                                        // VestingFunds
						fmt.Sprintf("%0.3f", util.ToFIL(res.VestingFunds)),               // VestingFundsFIL
						res.InitialPledge.String(),                                       // InitialPledge
						fmt.Sprintf("%0.3f", util.ToFIL(res.InitialPledge)),              // InitialPledgeFIL
						res.FeeDebt.String(),                                             // FeeDebt
						fmt.Sprintf("%0.3f", util.ToFIL(res.FeeDebt)),                    // FeeDebtFIL
						res.TerminationFeeFromSample.String(),                            // TerminationFee
						fmt.Sprintf("%0.3f", util.ToFIL(res.TerminationFeeFromSample)),   // TerminationFeeFIL
						res.AvgTerminationFeePerPledge.String(),                          // AvgTerminationFeePerPledge
						fmt.Sprintf("%0.3f", util.ToFIL(res.AvgTerminationFeePerPledge)), // AvgTerminationFeePerPledgeFIL
						fmt.Sprintf("%d", res.LiveSectors),                               // LiveSectors
						fmt.Sprintf("%d", res.FaultySectors),                             // FaultySectors
					}); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
					w.Flush()
					if err := w.Error(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().Bool("progress", true, "Output progress logs to stderr")
	collectCmd.Flags().Bool("debug", false, "Output debug logs to stderr")
	collectCmd.Flags().Bool("csv", true, "Output csv to stdout")
	collectCmd.Flags().String("save-sectors-dir", "", "If set, save CSV files with sector data for each miner in directory")
}
