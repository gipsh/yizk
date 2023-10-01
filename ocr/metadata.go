package ocr

import (
	"encoding/json"
	"math"
	"os"
	"strings"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gipsh/yizk/model"
	"go.uber.org/zap"
)

func (ocr *OcrService) generateMetadata(originalImage string, blocks []*visionpb.Block, pageId string) (*model.YizkPage, error) {

	page := model.YizkPage{
		Filename: originalImage,
		Id:       pageId,
	}

	for _, block := range blocks {
		if block.BlockType == visionpb.Block_TEXT {

			var sb strings.Builder
			pblock := model.YizkBlock{}
			//sb.WriteString("[Block]\n")

			nv := block.GetBoundingBox().GetVertices()
			maxX := getMaxX(nv)
			maxY := getMaxY(nv)
			minX := getMinX(nv)
			minY := getMinY(nv)

			ocr.log.Debug("Upper left point", zap.Int32("x", minX), zap.Int32("y", minY))
			ocr.log.Debug("Lower right point", zap.Int32("x", maxX), zap.Int32("y", maxY))

			pblock.UpperLeftPoint = model.YizkPoint{X: int(minX), Y: int(minY)}
			pblock.BottomRightPoint = model.YizkPoint{X: int(maxX), Y: int(maxY)}

			// p1 := model.YizkPoint{X: int(maxX), Y: int(maxY)}
			// p2 := model.YizkPoint{X: int(minX), Y: int(minY)}
			// pblock.Points = append(pblock.Points, p1, p2)

			langs := block.GetProperty().GetDetectedLanguages()
			if len(langs) > 0 {
				pblock.OriginLanguage = langs[0].GetLanguageCode()
			}

			for _, paragraph := range block.Paragraphs {
				//sb.WriteString("[Paragraph]\n")
				for _, word := range paragraph.Words {
					if ocr.isValidLang(word.Property.GetDetectedLanguages()) {
						for _, symbol := range word.Symbols {
							sb.WriteString(symbol.GetText())
						}
						sb.Write([]byte(" "))
					} else {
						ocr.log.Info("Skipping word", zap.Any("lang", word.Property.GetDetectedLanguages()))
					}
				}
			}
			sb.WriteString("\n")
			pblock.Text = sb.String()
			page.Blocks = append(page.Blocks, pblock)

		}
	}

	return &page, nil

}

func WriteMetadataFile(filename string, page *model.YizkPage) error {

	file, err := json.MarshalIndent(page, "", " ")
	if err != nil {
		return err
	}

	os.WriteFile(filename, file, 0644)

	return nil
}

func getMaxX(vertices []*visionpb.Vertex) int32 {
	var maxX int32
	for _, v := range vertices {
		if v.GetX() > maxX {
			maxX = v.GetX()
		}
	}
	return maxX
}

func getMaxY(vertices []*visionpb.Vertex) int32 {
	var maxY int32
	for _, v := range vertices {
		if v.GetY() > maxY {
			maxY = v.GetY()
		}
	}
	return maxY
}

func getMinX(vertices []*visionpb.Vertex) int32 {
	var minX int32
	minX = math.MaxInt32
	for _, v := range vertices {
		if v.GetX() < minX {
			minX = v.GetX()
		}
	}
	return minX
}

func getMinY(vertices []*visionpb.Vertex) int32 {
	var minY int32
	minY = math.MaxInt32
	for _, v := range vertices {
		if v.GetY() < minY {
			minY = v.GetY()
		}
	}
	return minY
}
