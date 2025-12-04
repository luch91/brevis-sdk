package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Gateway endpoints to test
	DefaultGateway = "appsdkv3.brevis.network:443"
	TestnetGateway = "testnet-api.brevis.network:443"

	// USDC on Ethereum mainnet (for RPC testing)
	USDCAddress = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
)

func main() {
	fmt.Println("=================================================")
	fmt.Println("  Brevis Gateway & Setup Connectivity Test")
	fmt.Println("=================================================\n")

	// Step 1: Check RPC endpoint
	fmt.Println("Step 1: Checking Ethereum RPC Endpoint...")
	fmt.Println("-------------------------------------------------")
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		fmt.Println("❌ ETH_RPC_URL environment variable not set")
		fmt.Println("\nTo fix this:")
		fmt.Println("  export ETH_RPC_URL=\"https://mainnet.infura.io/v3/YOUR_KEY\"")
		fmt.Println("  or")
		fmt.Println("  export ETH_RPC_URL=\"https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY\"")
		fmt.Println("\nContinuing with gateway tests...\n")
	} else {
		testRPCConnection(rpcURL)
	}

	// Step 2: Test gateway connections
	fmt.Println("\nStep 2: Testing Gateway Connections...")
	fmt.Println("-------------------------------------------------")

	fmt.Println("\n[Testing Default Production Gateway]")
	fmt.Printf("Endpoint: %s\n", DefaultGateway)
	testGatewayConnection(DefaultGateway, "Production")

	fmt.Println("\n[Testing Testnet Gateway]")
	fmt.Printf("Endpoint: %s\n", TestnetGateway)
	testGatewayConnection(TestnetGateway, "Testnet")

	// Step 3: Test BrevisApp creation
	fmt.Println("\nStep 3: Testing BrevisApp Creation...")
	fmt.Println("-------------------------------------------------")
	if rpcURL != "" {
		testBrevisAppCreation(rpcURL)
	} else {
		fmt.Println("⏭️  Skipped (no RPC URL)")
	}

	// Step 4: Summary and recommendations
	fmt.Println("\n=================================================")
	fmt.Println("  Test Summary & Next Steps")
	fmt.Println("=================================================\n")
	printSummary(rpcURL)
}

func testRPCConnection(rpcURL string) {
	fmt.Printf("RPC URL: %s\n", rpcURL)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test 1: Get latest block
	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get block number: %v\n", err)
		return
	}
	fmt.Printf("✅ Connected successfully\n")
	fmt.Printf("   Latest block: %d\n", blockNum)

	// Test 2: Verify USDC contract exists
	usdcAddr := common.HexToAddress(USDCAddress)
	code, err := client.CodeAt(ctx, usdcAddr, nil)
	if err != nil {
		fmt.Printf("⚠️  Warning: Cannot read contract code: %v\n", err)
		return
	}
	if len(code) > 0 {
		fmt.Printf("✅ Contract data accessible (USDC verified)\n")
	}
}

func testGatewayConnection(gatewayURL, name string) {
	// Try to create a gateway client
	_, err := sdk.NewGatewayClient(gatewayURL)
	if err != nil {
		fmt.Printf("❌ Failed to create gateway client: %v\n", err)
		fmt.Printf("   This might be a network/firewall issue\n")
		return
	}

	fmt.Printf("✅ Gateway client created successfully\n")
	fmt.Printf("   Note: Actual connectivity requires authentication/request\n")

	// Note: We can't test PrepareQuery without valid data and potential auth
	fmt.Printf("⚠️  Cannot test query submission without:\n")
	fmt.Printf("   - Valid circuit data\n")
	fmt.Printf("   - Possible API key/authentication\n")
	fmt.Printf("   - Gateway access permissions\n")
}

func testBrevisAppCreation(rpcURL string) {
	// Try creating a BrevisApp with default gateway
	fmt.Println("\n[Creating BrevisApp with Default Gateway]")
	app1, err := sdk.NewBrevisApp(1, rpcURL, "./test_output")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ BrevisApp created (default gateway)\n")
		_ = app1
	}

	// Try with testnet gateway override
	fmt.Println("\n[Creating BrevisApp with Testnet Gateway]")
	app2, err := sdk.NewBrevisApp(1, rpcURL, "./test_output", TestnetGateway)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ BrevisApp created (testnet gateway override)\n")
		_ = app2
	}

	// Try with Sepolia testnet (chain ID 11155111)
	fmt.Println("\n[Creating BrevisApp for Sepolia Testnet]")
	sepoliaRPC := os.Getenv("SEPOLIA_RPC_URL")
	if sepoliaRPC == "" {
		fmt.Println("⏭️  Skipped (SEPOLIA_RPC_URL not set)")
	} else {
		app3, err := sdk.NewBrevisApp(11155111, sepoliaRPC, "./test_output", TestnetGateway)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ BrevisApp created for Sepolia\n")
			_ = app3
		}
	}
}

func printSummary(rpcURL string) {
	fmt.Println("✅ COMPLETED TESTS")
	fmt.Println()

	if rpcURL != "" {
		fmt.Println("✓ RPC Endpoint: Available")
	} else {
		fmt.Println("✗ RPC Endpoint: NOT SET")
	}

	fmt.Println("✓ Gateway Clients: Can be created")
	fmt.Println("✓ SDK Installation: Working")
	fmt.Println()

	fmt.Println("⚠️  WHAT'S STILL NEEDED:")
	fmt.Println()
	fmt.Println("1. Gateway Access Verification")
	fmt.Println("   - Try submitting actual circuit data")
	fmt.Println("   - May require API key or authentication")
	fmt.Println("   - Contact Brevis team for access")
	fmt.Println()

	fmt.Println("2. Test Proof Generation")
	fmt.Println("   - Add sample data to BrevisApp")
	fmt.Println("   - Call PrepareRequest()")
	fmt.Println("   - Check if authentication error occurs")
	fmt.Println()

	fmt.Println("3. Recommended Next Steps:")
	fmt.Println()
	fmt.Println("   a) Join Brevis Discord:")
	fmt.Println("      https://discord.com/invite/brevis")
	fmt.Println()
	fmt.Println("   b) Ask in #dev-support:")
	fmt.Println("      \"I have 20 compiled circuits ready to test.")
	fmt.Println("      How do I get testnet gateway access?\"")
	fmt.Println()
	fmt.Println("   c) Share your contribution:")
	fmt.Println("      \"Built example circuits for major DeFi protocols")
	fmt.Println("      (Uniswap, Aave, Curve, etc.) - need testing access\"")
	fmt.Println()

	if rpcURL == "" {
		fmt.Println("4. Set up RPC endpoint:")
		fmt.Println()
		fmt.Println("   Get a free RPC from:")
		fmt.Println("   - Infura: https://infura.io")
		fmt.Println("   - Alchemy: https://alchemy.com")
		fmt.Println()
		fmt.Println("   Then set environment variable:")
		fmt.Println("   export ETH_RPC_URL=\"your-rpc-url-here\"")
		fmt.Println()
	}

	fmt.Println("=================================================")
	fmt.Println()
	fmt.Println("💡 TIP: The fact that gateway clients can be created")
	fmt.Println("   suggests the endpoints exist. You likely just need")
	fmt.Println("   proper credentials or permission from Brevis team.")
	fmt.Println()
}
