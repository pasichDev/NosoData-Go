package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Friends-Of-Noso/NosoData-Go/legacy"
	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

const (
	cTestDataFolder  = "testdata"
	cBlockFilename   = "100000.blk"
	cWalletFilename  = "wallet.pkw"
	cSummaryFilename = "sumary.psk"
	cGVTFilename     = "gvts.psk"
	cPSOFilename     = "psos.dat"
)

var (
	block   legacy.LegacyBlock
	wallet  legacy.LegacyWallet
	summary legacy.LegacySummary
	gvts    legacy.LegacyGVT
	psos    legacy.LegacyPSO
)

func displayBlock(jsonOutput bool) {
	fmt.Printf("\n%s\n", "== Block ==")
	err := block.ReadFromFile(filepath.Join(cTestDataFolder, cBlockFilename))
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	if jsonOutput {
		fmt.Println(block.AsJSON())
	} else {

		fmt.Println("Number:           ", block.Number)
		fmt.Printf("HASH:              '%s'\n", block.HASH)
		fmt.Println("Time Start:       ", time.Unix(block.TimeStart, 0))
		fmt.Println("Time End:         ", time.Unix(block.TimeEnd, 0))
		fmt.Println("Time Total:       ", block.TimeTotal, "seconds")
		fmt.Println("Time Last 20:     ", block.TimeLast20, "seconds")
		fmt.Println("Transaction Count:", block.TransactionsCount)
		fmt.Println("Difficulty:       ", block.Difficulty)
		fmt.Printf("Target Hash:      '%s'\n", block.TargetHash.GetString())
		fmt.Printf("Solution:         '%s'\n", block.Solution.GetString())
		fmt.Printf("Last Block Hash:  '%s'\n", block.LastBlockHash.GetString())
		fmt.Printf("Miner:            '%s'\n", block.Miner.GetString())
		fmt.Println("Fee:             ", utils.ToNoso(block.Fee))
		fmt.Println("Reward:          ", utils.ToNoso(block.Reward))

		if block.TransactionsCount > 0 {
			fmt.Printf("Transactions(%d):\n", block.TransactionsCount)
			var n int32
			for n = 0; n < block.TransactionsCount; n++ {
				fmt.Printf("  OrderID: '%s'\n", block.Transactions[n].OrderID.GetString())
				fmt.Printf("      TransferID:     '%s'\n", block.Transactions[n].TransferID.GetString())
				fmt.Println("      Block:         ", block.Transactions[n].Block)
				fmt.Println("      Order lines:   ", block.Transactions[n].OrderLinesCount)
				fmt.Printf("      Order type:     '%s'\n", block.Transactions[n].OrderType.GetString())
				fmt.Println("      Timestamp:     ", time.Unix(block.Transactions[n].TimeStamp, 0))
				fmt.Printf("      Reference:      '%s'\n", block.Transactions[n].Reference.GetString())
				fmt.Println("      Transfer Index:", block.Transactions[n].TransferIndex)
				fmt.Printf("      Sender:         '%s'\n", block.Transactions[n].Sender.GetString())
				fmt.Printf("      Address:        '%s'\n", block.Transactions[n].Address.GetString())
				fmt.Printf("      Receiver:       '%s'\n", block.Transactions[n].Receiver.GetString())
				fmt.Println("      Fee:           ", utils.ToNoso(block.Transactions[n].AmountFee))
				fmt.Println("      Value:         ", utils.ToNoso(block.Transactions[n].AmountTransfer))
				fmt.Printf("      Signature:      '%s'\n", block.Transactions[n].Signature.GetString())
			}
		} else {
			fmt.Println("No transactions")
		}

		if block.ProofOfStakeRewardCount > 0 {
			fmt.Printf("PoS rewards(%d):\n", block.ProofOfStakeRewardCount)
			fmt.Println("  Amount:", utils.ToNoso(block.ProofOfStakeRewardAmount))
			var n int32
			for n = 0; n < block.ProofOfStakeRewardCount; n++ {
				fmt.Printf("  Address: '%s'\n", block.ProofOfStakeRewardAddresses[n].GetString())
			}
		} else {
			fmt.Println("No PoS rewards")
		}

		if block.MasterNodeRewardCount > 0 {
			fmt.Printf("MN rewards(%d):\n", block.MasterNodeRewardCount)
			fmt.Println("  Amount:", utils.ToNoso(block.MasterNodeRewardAmount))
			var n int32
			for n = 0; n < block.MasterNodeRewardCount; n++ {
				fmt.Printf("  Address: '%s'\n", block.MasterNodeRewardAddresses[n].GetString())
			}
		} else {
			fmt.Println("No MN rewards")
		}
	}
}

func displayWallet(jsonOutput bool) {
	// Wallet
	fmt.Printf("\n%s\n", "== Wallet ==")
	err := wallet.ReadFromFile(filepath.Join(cTestDataFolder, cWalletFilename))
	if err != nil {
		fmt.Println("error reading wallet:", err)
		return
	}
	if jsonOutput {
		fmt.Println(wallet.AsJSON())
	} else {
		for i, a := range wallet.Accounts {
			fmt.Println("Position:", i)
			fmt.Printf("    Hash: '%s'\n", a.Hash.GetString())
			fmt.Printf("    Custom:         '%s'\n", a.Custom.GetString())
			fmt.Printf("    Pub key:        '%s'\n", a.PublicKey.GetString())
			fmt.Printf("    Priv key:       '%s'\n", a.PrivateKey.GetString())
			fmt.Println("    Balance:       ", utils.ToNoso(a.Balance))
			fmt.Println("    Pending:       ", utils.ToNoso(a.Pending))
			fmt.Println("    Score:         ", utils.ToNoso(a.Score))
			fmt.Println("    Last Operation:", utils.ToNoso(a.LastOperation))
		}
	}
}

func displaySummary(jsonOutput bool) {
	// Summary
	fmt.Printf("\n%s\n", "== Summary ==")
	err := summary.ReadFromFile(filepath.Join(cTestDataFolder, cSummaryFilename))
	if err != nil {
		fmt.Println("error reading summary:", err)
		return
	}
	if jsonOutput {
		fmt.Println(summary.AsJSON())
	} else {
		for i, a := range summary.Accounts {
			fmt.Println("Position:", i)
			fmt.Printf("    Hash:           '%s'\n", a.Hash.GetString())
			fmt.Printf("    Custom:         '%s'\n", a.Custom.GetString())
			fmt.Println("    Balance:       ", utils.ToNoso(a.Balance))
			fmt.Println("    Score:         ", utils.ToNoso(a.Score))
			fmt.Println("    Last Operation:", utils.ToNoso(a.LastOperation))
		}
	}

}

func displayGVT(jsonOutput bool) {
	// GVT
	fmt.Printf("\n%s\n", "== GVT ==")
	err := gvts.ReadFromFile(filepath.Join(cTestDataFolder, cGVTFilename))
	if err != nil {
		fmt.Println("error reading GVT:", err)
		return
	}
	if jsonOutput {
		fmt.Println(gvts.AsJSON())
	} else {
		for i, e := range gvts.Entries {
			fmt.Println("Position:", i)
			fmt.Printf("    Number:  '%s'\n", e.Number.GetString())
			fmt.Printf("    Owner:   '%s'\n", e.Owner.GetString())
			fmt.Printf("    Hash:    '%s'\n", e.Hash.GetString())
			fmt.Println("    Control:", e.Control)
		}
	}
}

func displayPSO(jsonOutput bool) {
	// PSO
	fmt.Printf("\n%s\n", "== PSO ==")
	err := psos.ReadFromFile(filepath.Join(cTestDataFolder, cPSOFilename))
	if err != nil {
		fmt.Println("error reading PSO:", err)
		return
	}
	if jsonOutput {
		fmt.Println(psos.AsJSON())
	} else {
		fmt.Println("Block:", psos.Block)
		fmt.Printf("  MN Locks(%d):\n", psos.MNLockCount)
		for i, mli := range psos.MNLocks {
			fmt.Println("  Position:", i)
			fmt.Printf("      Address: '%s'\n", mli.Address.GetString())
			fmt.Println("       Expire:", mli.Expire, "seconds")
		}
		fmt.Printf("  PSO Count(%d):\n", psos.PSOCount)
	}

}

func main() {

	displayPSO(true)

	displayGVT(true)

	displaySummary(true)

	displayWallet(true)

	displayBlock(true)
}
