package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGettingSiteTemplate(t *testing.T) {
	expected := `
<!doctype html>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="generator" content="Jordans Order Test">
	<meta name="viewport" content="width=device-width, minimum-scale=1, initial-scale=1, user-scalable=yes">

	<title>Jordans Order Test</title>
	<meta name="description" content="Jordans Order Test">
	<meta name="author" content="https://github.com/jordanfinners">

	<meta property="og:title" content="Jordans Order Test">
	<meta property="og:description" content="Jordans Order Test">
	<meta property="og:image"
		content="https://avatars2.githubusercontent.com/u/17813098?s=460&u=f8f61170c39933eff8aaf52f87bf6939ecdee7a6&v=4">
</head>

<body>
	<form name="order" method="post" action="{{ .ActionURL }}">
		<label>How many items do you wish to order?
			<input type="number" name="items" placeholder="e.g. 501" min="1" required>
		</label>
		<button type="submit">Submit Order</button>
		<button type="reset">Start Over</button>
	</form>
</body>

</html>`

	template := getSiteTemplate()

	require.Equal(t, expected, template)
}
