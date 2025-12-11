package tools

import (
	"context"
	"encoding/json"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// payramIntroContent is returned by the payram_intro tool.
const payramIntroContent = `Welcome to PayRam

The complete self-hosted payments stack for global commerce.
Accept payments globally, monetize anything in minutes, with no middlemen, censorship-resistant settlement, and full custody, data, and control on your own infrastructure.

What is PayRam?
PayRam is a self-hosted PayFi platform for stablecoin and cryptocurrency payments. It lets you accept and settle payments directly onchain, with no middlemen, no custody risk, and full control over your funds and data. Built for financial liberalization, PayRam is censorship-resistant, programmable, and designed to help anyone run global commerce on infrastructure they own without any middlemen.

Get started with PayRam
- Deployment Guide: https://docs.payram.com/deployment-guide/introduction
- Onboarding Guide: https://docs.payram.com/onboarding-guide/introduction
- PayRam Features: https://docs.payram.com/features/payment-links

Need help?
- Community Support: https://x.com/PayRamApp
- Contact Support: https://payram.short.gy/payram-gitbook-contact
`

// payramIntroTool implements the PayRam intro tool.
type payramIntroTool struct{}

// PayramIntro constructs the PayRam intro tool instance.
func PayramIntro() *payramIntroTool {
	return &payramIntroTool{}
}

func (t *payramIntroTool) Descriptor() protocol.ToolDescriptor {
	return protocol.ToolDescriptor{
		Name:        "payram_intro",
		Description: "Overview of PayRam and helpful links.",
	}
}

func (t *payramIntroTool) Invoke(ctx context.Context, raw json.RawMessage) (protocol.CallResult, *protocol.ResponseError) {
	return protocol.CallResult{
		Content: []protocol.ContentPart{{Type: "text", Text: payramIntroContent}},
	}, nil
}
