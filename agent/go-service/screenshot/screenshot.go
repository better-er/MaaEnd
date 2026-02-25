package screenshot

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	maa "github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/rs/zerolog/log"
)

// ScreenShot 将当前截图保存为 PNG 到 debug 目录，可用于调试。
// custom_action_param 为 JSON，可选字段 "type" 作为文件名前缀。
type ScreenShot struct{}

var _ maa.CustomActionRunner = (*ScreenShot)(nil)

// Run 实现 maa.CustomActionRunner：截屏并保存为 PNG。
func (a *ScreenShot) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	var params struct {
		Type string `json:"type"`
	}
	if arg.CustomActionParam != "" {
		_ = json.Unmarshal([]byte(arg.CustomActionParam), &params)
	}
	typePrefix := strings.TrimSpace(params.Type)
	if typePrefix != "" {
		typePrefix = typePrefix + "_"
	}

	ctrl := ctx.GetTasker().GetController()
	ctrl.PostScreencap().Wait()
	img, err := ctrl.CacheImage()
	if err != nil {
		log.Error().Err(err).Msg("[ScreenShot] 截图失败")
		return false
	}
	if img == nil {
		log.Error().Msg("[ScreenShot] 截图为空")
		return false
	}

	debugDir, err := filepath.Abs("debug")
	if err != nil {
		log.Error().Err(err).Msg("[ScreenShot] 解析 debug 目录失败")
		return false
	}
	if err := os.MkdirAll(debugDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", debugDir).Msg("[ScreenShot] 创建 debug 目录失败")
		return false
	}

	// 清理 3 天前的 PNG 文件
	threeDaysAgo := time.Now().Add(-3 * 24 * time.Hour)
	entries, _ := os.ReadDir(debugDir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(e.Name()), ".png") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(threeDaysAgo) {
			p := filepath.Join(debugDir, e.Name())
			if err := os.Remove(p); err != nil {
				log.Debug().Err(err).Str("path", p).Msg("[ScreenShot] 清理旧文件失败")
			}
		}
	}

	fileName := typePrefix + time.Now().Format("2006-01-02_15-04-05") + ".png"
	debugPath := filepath.Join(debugDir, fileName)

	// 若 CacheImage 返回的是非 *image.RGBA，转为 RGBA 以便编码
	toEncode := image.Image(img)
	if _, ok := img.(*image.RGBA); !ok {
		b := img.Bounds()
		rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
		toEncode = rgba
	}

	f, err := os.Create(debugPath)
	if err != nil {
		log.Error().Err(err).Str("path", debugPath).Msg("[ScreenShot] 创建文件失败")
		return false
	}
	defer f.Close()
	if err := png.Encode(f, toEncode); err != nil {
		log.Error().Err(err).Str("path", debugPath).Msg("[ScreenShot] 写入 PNG 失败")
		return false
	}
	log.Info().Str("path", debugPath).Msg("[ScreenShot] 已保存截图")
	return true
}
