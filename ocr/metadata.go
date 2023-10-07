package ocr

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gipsh/yizk/model"
	"go.uber.org/zap"
)

func (ocr *OcrService) generateMetadata(originalImage string, blocks []*visionpb.Block, pageId string, pageNum int) (*model.YizkPage, error) {

	page := model.YizkPage{
		Filename: originalImage,
		Id:       pageId,
		Order:    pageNum,
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

			langs := block.GetProperty().GetDetectedLanguages()
			if len(langs) > 0 {
				pblock.OriginLanguage = langs[0].GetLanguageCode()
			}

			// TODO: add language/confidence histogram to block
			for _, paragraph := range block.Paragraphs {
				for _, word := range paragraph.Words {
					if ocr.isValidLang(word.Property.GetDetectedLanguages()) {
						for _, symbol := range word.Symbols {
							sb.WriteString(symbol.GetText())
						}
						sb.Write([]byte(" "))
					} else {
						ocr.log.Debug("Skipping word", zap.Any("lang", word.Property.GetDetectedLanguages()), zap.Any("text", word.Symbols))
					}
				}
			}

			// remove double quotes from text (they break the json) and the translation
			// WARNING: seems the double yud is detected as `"`. This is a hack, need to find a better way to do this
			pblock.Text = strings.ReplaceAll(sb.String(), `"`, ``)
			//pblock.Text = sb.String()

			if len(pblock.Text) > 0 && pblock.Text != " " && pblock.Text != "\n" {
				ip, pn := ocr.isPageNumber(pblock.Text)
				if ip {
					ocr.log.Info("Found page number", zap.String("page", pn))
					page.PageNumber = pn
				} else {
					page.Blocks = append(page.Blocks, pblock)
				}
			} else {
				ocr.log.Info("Skipping block", zap.Any("block", pblock), zap.Any("text", pblock.Text))
			}
		}
	}

	return &page, nil

}

func (ocr *OcrService) isPageNumber(blockText string) (bool, string) {
	text := strings.ReplaceAll(blockText, " ", "")
	text = strings.ReplaceAll(text, "\n", "")

	_, err := strconv.Atoi(text)
	return err == nil, text
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
