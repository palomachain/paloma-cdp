package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type subscribeInput struct {
	// SymbolName string `path:"name" required:"true" pattern:"^(DEX|BONDING):\\S{3,44}-[a-z0-9]{6}/\\S{3,44}-[a-z0-9]{6}$" description:"Full symbol name, including leading exchange name." example:"BONDING:UPUSD-19nr9t/MTT.0-1wff3z"`
	SymbolName string `path:"name" required:"true" `
	// Resolution string `query:"resolution" required:"true" enum:"1S,2S,5S,1,2,5,60,120,300,1D,2D,1W,2W,1M,2M,3M" description:"Resolution of the symbol"`
	r *http.Request
	_ struct{} `query:"_" cookie:"_" additionalProperties:"false"`
}

func (s *subscribeInput) SetRequest(r *http.Request) {
	s.r = r
}

type subscribeOutput struct {
	w http.ResponseWriter
}

func (s *subscribeOutput) SetWriter(w io.Writer) {
	s.w = w.(http.ResponseWriter)
}

func SubscribeInteractor(db *persistence.Database) usecase.IOInteractor {
	u := usecase.NewInteractor(func(ctx context.Context, input subscribeInput, output *subscribeOutput) error {
		_, err := url.QueryUnescape(input.SymbolName)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}

		fmt.Println("r", input.r, "o", output.w)
		c, err := websocket.Accept(output.w, input.r, nil)
		if err != nil {
			// TODO: NOPE
			panic(err)
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			// TODO: NOPE
			panic(err)
		}

		fmt.Println("v", v)

		c.Close(websocket.StatusNormalClosure, "")

		return nil
	})
	u.SetTitle("Get Bars 2 lol")
	u.SetDescription("Returns a set of bars for a given symbol.")
	u.SetExpectedErrors(status.InvalidArgument, status.NotFound)
	u.SetTags(cTagAdvancedCharts)

	return u.IOInteractor
}
