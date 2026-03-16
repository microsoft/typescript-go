package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatTemplateLiteralInMalformedPropertySignatureNoCrash1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `// @Filename: /resource-card.tsx
interface Props {
  iconOnly?: boolean


const ResourceCard: React.FC<Props> = (props) => {
  return (
    <IoLayersOutline
      className={` + "`" + `${match ? 'text-primary-foreground' : 'text-foreground'}` + "`" + `}
    />
  )
}

export default ResourceCard
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `interface Props {
    iconOnly?: boolean


const ResourceCard: React.FC<Props> = (props) => {
    return (
        <IoLayersOutline
            className={`+"`"+`${match ? 'text-primary-foreground' : 'text-foreground'}`+"`"+`}
        />
    )
}

export default ResourceCard
`)
}
