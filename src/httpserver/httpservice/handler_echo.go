package httpservice

import (
	"context"
	mylog "endocode/src/log"
	"fmt"
	myctx "github.com/romapres2010/httpserver/ctx"
	"net/http"
)

// EchoHandler handle echo page with request header and body
func (s *Service) EchoHandler(w http.ResponseWriter, r *http.Request) {
	mylog.PrintfDebugMsg("START   ==================================================================================")

	// Запускаем обработчик, возврат ошибки игнорируем
	_ = s.process("POST", w, r, func(ctx context.Context, requestBuf []byte, buf []byte) ([]byte, Header, int, error) {
		reqID := myctx.FromContextRequestID(ctx) // RequestID передается через context

		mylog.PrintfDebugMsg("START: reqID", reqID)

		// формируем ответ
		header := Header{}
		header["Errcode"] = "0"
		header["RequestID"] = fmt.Sprintf("%v", reqID)

		// Считаем параметры из заголовка сообщения и перенесем их в ответный заголовок
		for key := range r.Header {
			header[key] = r.Header.Get(key)
		}

		mylog.PrintfDebugMsg("SUCCESS", reqID)

		// входной буфер возвращаем в качестве выходного
		return requestBuf, header, http.StatusOK, nil
	})

	mylog.PrintfDebugMsg("SUCCESS ==================================================================================")
}
