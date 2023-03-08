package pl

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("pl", parseCaddyfile)
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	ShutdownFile string `json:"shutdown_file"`
}

// 注册插件
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.pl",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	// 处理插件业务
	go m.signListen()
	return nil
}

// 验证数据可以不需要
func (m *Middleware) Validate() error {
	// if m.ShutdownFile == "" {
	// 	return errors.New("the text is must!!!")
	// }
	return nil
}

// 修改数据,不修改直接返回
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// w.Write([]byte(m.ShutdownFile))
	// return next.ServeHTTP(w, r)
	// 如果文件存在,直接500
	if FileExist(m.ShutdownFile) {
		w.WriteHeader(500)
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 验证解析和保存
func (m *Middleware) UnmarshalCaddyfile(h *caddyfile.Dispenser) error {
	// 解析参数并赋值
	for h.Next() {
		for h.NextBlock(0) {
			opt := h.Val()
			switch opt {
			case "shutdown_file":
				// 校验和赋值
				if !h.Args(&m.ShutdownFile) {
					return h.ArgErr()
				}
			}
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)

// 信号监听
func (m *Middleware) signListen() {
	c := make(chan os.Signal)
	signal.Notify(c)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				m.filelog("服务即将停止...")
			default:
				m.filelog("other signal")
			}
		}
	}()
}

// 写文件
func (m *Middleware) filelog(content string) {
	filename := m.ShutdownFile
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		// 创建文件失败处理
	} else {
		_, err = f.Write([]byte(content))
		if err != nil {
			// 写入失败处理
		}
	}
}
