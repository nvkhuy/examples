package tests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"os"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/html"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestPost_GenerateSlug(t *testing.T) {
	var app = initApp("local")
	err := repo.NewPostRepo(app.DB).GenerateSlug()

	assert.NoError(t, err)
}

func TestPost_PaginatePost(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewPostRepo(app.DB).PaginatePost(repo.PaginatePostParams{
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 12,
		},
		Statuses: []enums.PostStatus{
			enums.PostStatusPublished,
		},
		//Language: enums.LanguageCodeVietnam,
	})

	helper.PrintJSON(result)
}

func TestPost_GetPost(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewPostRepo(app.DB).GetPost(repo.GetPostParams{
		Slug: "content-test-update",
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestPost_GetPostStats(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewPostRepo(app.DB).GetPostStats(repo.GetPostStatsParams{})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestPost_CreatePost(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewPostRepo(app.DB).UpdatePost(models.PostUpdateForm{
		PostID:         "cl8tuprq2n6gtpmt2gjg",
		PostCreateForm: models.PostCreateForm{},
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestPost_DownloadContent(t *testing.T) {
	var app = initApp("prod")
	var result *models.Post
	var err error
	result, err = repo.NewPostRepo(app.DB).GetPost(repo.GetPostParams{
		Slug: "mastering-apparel-sourcing-and-manufacturing-in-vietnam-inflows-comprehensive-guide",
	})

	filePath := fmt.Sprintf("post-en-%s.html", result.ID)

	err = os.WriteFile(filePath, []byte(result.Content), 0777)
	if err != nil {
		fmt.Println("Error writing HTML to file:", err)
		return
	}

	fmt.Printf("HTML string saved to %s\n", filePath)
}

func TestPost_UploadCompressedHTML(t *testing.T) {
	var app = initApp("prod")
	var result *models.Post
	var err error
	result, err = repo.NewPostRepo(app.DB).GetPost(repo.GetPostParams{
		Slug: "recap-brand-ta-ra-tay-18-11-a-leap-forward-for-local-fashion-brands-in-vietnam",
	})

	// extract start
	htmlFile := fmt.Sprintf("uploads/content/post-en-%s", result.ID)
	prefix := fmt.Sprintf("uploads/content/post-en-media-%s", result.ID)
	bucket := app.Config.AWSS3StorageBucket
	storageURL := app.Config.StorageURL

	extractedHTML, _, err := html.NewHTMLCompress(htmlFile, prefix, result.Content, bucket, storageURL, app.S3Client).Compress()
	if err != nil {
		return
	}
	log.Println(len(extractedHTML))

	// end
	filePath := fmt.Sprintf("post-en-%s.html", result.ID)
	err = os.WriteFile(filePath, []byte(extractedHTML), 0777)
	if err != nil {
		fmt.Println("Error writing HTML to file:", err)
		return
	}

	fmt.Printf("HTML string saved to %s\n", filePath)
}

func TestPost_UpdateContentURL(t *testing.T) {
	var app = initApp("local")
	var err = repo.NewPostRepo(app.DB).UpdateContentURL("clkqdenskmp3nd6qgct0", enums.LanguageCodeEnglish, "https://upload/post-1.html")
	assert.NoError(t, err)
}

func TestPost_MoveImageToS3(t *testing.T) {
	var app = initApp("prod")
	var posts []models.Post
	var err error
	err = app.DB.Select("id").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return
	}
	count := float64(0)
	total := float64(len(posts))
	for _, p := range posts {
		if err = repo.NewPostRepo(app.DB).MoveImageToS3(p.ID); err != nil {
			log.Println(err)
		}
		count += 1
		progress := fmt.Sprintf("%f", (count/total)*100) + "%"
		fmt.Printf("BLog Post %s Finish - Progress: %s\n", p.ID, progress)
	}
	assert.NoError(t, err)
}

func TestPost_PaginateBlog(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewBlogCategoryRepo(app.DB).PaginateBlogCategory(repo.PaginateBlogCategoryParams{
		TotalPost: aws.Int(0),
	})
	helper.PrintJSON(result)
}
