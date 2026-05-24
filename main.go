package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zserge/lorca"
)

func main() {
	// Chromeのセキュリティ制限回避パラメータを指定
	ui, err := lorca.New("", "", 600, 650, "--remote-allow-origins=*")
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// Go側での液性判定ロジック
	ui.Bind("goJudgePH", func(ph float64) string {
		var result string
		switch {
		case ph < 0.0 || ph > 14.0:
			result = "範囲外エラー"
		case ph < 7.0:
			result = fmt.Sprintf("酸性 (pH %.1f)", ph)
		case ph == 7.0:
			result = fmt.Sprintf("中性 (pH %.1f)", ph)
		default:
			result = fmt.Sprintf("アルカリ性 (pH %.1f)", ph)
		}

		response := map[string]string{"liquidType": result}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	})

	// ドル記号（$）を完全に排除したHTML/CSS/JS
	htmlContent := `
	<!DOCTYPE html>
	<html lang="ja">
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
				padding: 20px;
				background-color: #f5f5f7;
				color: #333;
				user-select: none;
			}
			.container {
				max-width: 500px;
				margin: 0 auto;
				background: white;
				padding: 25px;
				border-radius: 12px;
				box-shadow: 0 4px 6px rgba(0,0,0,0.1);
			}
			h3 { margin-top: 0; color: #111; text-align: center; }
			.slider-box { margin: 20px 0; text-align: center; }
			input[type="range"] { width: 100%; margin-top: 10px; }
			.status { font-weight: bold; font-size: 1.1em; text-align: center; color: #0071e3; min-height: 24px; }
			.indicator-box {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 15px;
				margin-top: 25px;
			}
			.card {
				border: 1px solid #ddd;
				border-radius: 8px;
				padding: 15px;
				text-align: center;
				background: #fafafa;
			}
			.label { font-size: 0.9em; color: #666; margin-bottom: 8px; }
			.litmus {
				width: 30px;
				height: 100px;
				margin: 0 auto;
				border-radius: 4px;
				transition: background-color 0.3s;
				box-shadow: inset 0 0 5px rgba(0,0,0,0.2);
			}
			.tube {
				width: 40px;
				height: 100px;
				margin: 0 auto;
				border: 2px solid #bbb;
				border-bottom-left-radius: 20px;
				border-bottom-right-radius: 20px;
				position: relative;
				background: #eee;
				overflow: hidden;
			}
			.liquid {
				position: absolute;
				bottom: 0;
				left: 0;
				right: 0;
				height: 80%;
				transition: background-color 0.3s;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h3>pHシミュレーター (Lorca 安定版)</h3>
			<div class="slider-box">
				<label>pH値を調整してください</label>
				<input type="range" id="phSlider" min="0" max="14" step="0.1" value="7.0">
				<div style="display:flex; justify-content:space-between; font-size:0.8em; color:#888; padding:0 5px;">
					<span>0 (強酸)</span><span>7 (中性)</span><span>14 (強アルカリ)</span>
				</div>
			</div>
			<div class="status" id="statusText">中性 (pH 7.0)</div>
			<div class="indicator-box">
				<div class="card">
					<div class="label">赤色リトマス紙</div>
					<div id="redLitmus" class="litmus" style="background-color: #e57373;"></div>
				</div>
				<div class="card">
					<div class="label">青色リトマス紙</div>
					<div id="blueLitmus" class="litmus" style="background-color: #64b5f6;"></div>
				</div>
				<div class="card">
					<div class="label">BTB溶液</div>
					<div class="tube"><div id="btbLiquid" class="liquid" style="background-color: #81c784;"></div></div>
				</div>
				<div class="card">
					<div class="label">フェノールフタレイン</div>
					<div class="tube"><div id="ppLiquid" class="liquid" style="background-color: rgba(255,255,255,0.5);"></div></div>
				</div>
			</div>
		</div>

		<script>
			const slider = document.getElementById('phSlider');
			const statusText = document.getElementById('statusText');
			const redLitmus = document.getElementById('redLitmus');
			const blueLitmus = document.getElementById('blueLitmus');
			const btbLiquid = document.getElementById('btbLiquid');
			const ppLiquid = document.getElementById('ppLiquid');

			async function updateIndicators(ph) {
				const goResponseJson = await window.goJudgePH(parseFloat(ph));
				const data = JSON.parse(goResponseJson);
				statusText.innerText = data.liquidType;

				if (ph >= 8.0) {
					redLitmus.style.backgroundColor = '#42a5f5'; 
				} else {
					redLitmus.style.backgroundColor = '#e57373'; 
				}

				if (ph <= 5.0) {
					blueLitmus.style.backgroundColor = '#e57373';
				} else {
					blueLitmus.style.backgroundColor = '#64b5f6';
				}

				if (ph < 6.0) {
					btbLiquid.style.backgroundColor = '#fff176'; 
				} else if (ph <= 7.6) {
					btbLiquid.style.backgroundColor = '#81c784'; 
				} else {
					btbLiquid.style.backgroundColor = '#1e88e5'; 
				}

				if (ph < 8.3) {
					ppLiquid.style.backgroundColor = 'rgba(255,255,255,0.5)'; 
				} else {
					const intensity = (ph - 8.3) / (14 - 8.3);
					const opacity = 0.4 + intensity * 0.6;
					// ドル記号（$）を使わず、配列を結合させて安全に文字列化する
					ppLiquid.style.backgroundColor = ['rgba(240, 98, 146, ', opacity, ')'].join('');
				}
			}

			slider.addEventListener('input', (e) => {
				updateIndicators(e.target.value);
			});
		</script>
	</body>
	</html>
	`

	// 文字列をBase64形式にエンコードして安全に読み込ませる
	opaqueData := base64.StdEncoding.EncodeToString([]byte(htmlContent))
	ui.Load("data:text/html;base64," + opaqueData)

	<-ui.Done()
}