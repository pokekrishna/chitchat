package httpserver

import (
	"bytes"
	"github.com/pokekrishna/chitchat/internal/data"
	"io"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestGenerateHTML(t *testing.T){
	t.Run("Invalid filename should return error and return empty response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := "testdata"
		dummyFileName := "file"+strconv.FormatInt(time.Now().UnixNano(), 10)

		defer os.Chdir(ChangeDirForTest())
		err := generateHTML(w, data, dummyFileName)
		if err == nil{
			t.Errorf("Expected error but didnt get.")
		}

		bodyContent, err := w.Body.ReadByte()
		if err != io.EOF {
			t.Errorf("Expected blank body, Got Body: %s", string(bodyContent))
			t.Errorf("Expected error EOF, got: %s", err)
		}
	})

	t.Run("When inputs are valid", func(t *testing.T) {
		time1 := time.Now()
		time2 := time.Now()
		threads := []*data.Thread{
			&data.Thread{
				Id:        99,
				Uuid:      "e692d3d4-cb61-11eb-b8bc-0242ac130003",
				Topic:     "Cats",
				UserId:    102,
				CreatedAt: time1,
			},
			&data.Thread{
				Id:        70,
				Uuid:      "9992d3d4-cb61-11eb-b8bc-0242ac130111",
				Topic:     "Dogs",
				UserId:    111,
				CreatedAt: time2,
			},
		}
		expectedGeneratedHtml := `<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8">
				<meta http-equiv="X-UA-Compatible" content="IE=9">
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<title>ChitChat</title>
				<link href="/static/css/bootstrap.min.css" rel="stylesheet">
				<link href="/static/css/font-awesome.min.css" rel="stylesheet">
			</head>
			<body>
		
		<div class="navbar navbar-default navbar-static-top" role="navigation">
			<div class="container">
				<div class="navbar-header">
					<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target=".navbar-collapse">
						<span class="sr-only">Toggle navigation</span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="/">
						<i class="fa fa-comments-o"></i>
						ChitChat
					</a>
				</div>
				<div class="navbar-collapse collapse">
					<ul class="nav navbar-nav">
						<li><a href="/">Home</a></li>
					</ul>
					<ul class="nav navbar-nav navbar-right">
						<li><a href="/login">Login</a></li>
					</ul>
				</div>
			</div>
		</div>
		
				<div class="container">
		
		<p class="lead">
			<a href="/thread/new">Start a thread</a> or join one below!
		</p>
		
			<div class="panel panel-default">
				<div class="panel-heading">
					<span class="lead"> <i class="fa fa-comment-o"></i> Cats</span>
				</div>
				<div class="panel-body">
					Started by 102 - {{.Time1}} posts.
					<div class="pull-right">
						<a href="/thread/read?id=e692d3d4-cb61-11eb-b8bc-0242ac130003">Read more</a>
					</div>
				</div>
			</div>
		
			<div class="panel panel-default">
				<div class="panel-heading">
					<span class="lead"> <i class="fa fa-comment-o"></i> Dogs</span>
				</div>
				<div class="panel-body">
					Started by 111 - {{.Time2}} posts.
					<div class="pull-right">
						<a href="/thread/read?id=9992d3d4-cb61-11eb-b8bc-0242ac130111">Read more</a>
					</div>
				</div>
			</div>
		
		
				</div> <!-- /container -->
				<script src="/static/js/jquery-2.1.1.min.js"></script>
				<script src="/static/js/bootstrap.min.js"></script>
			</body>
		</html>`
		type Timestamps struct{
			Time1, Time2 time.Time
		}
		ts := Timestamps{
			Time1: time1,
			Time2: time2,
		}

		var b bytes.Buffer
		tpl, err := template.New("expectedGeneratedHTML").Parse(expectedGeneratedHtml)
		if err != nil {
			t.Errorf("Cannot Parse %s Template", tpl.Name())
		}
		if err := tpl.Execute(&b, ts); err != nil {
			t.Errorf("Cannot Execute %s Template", tpl.Name())
		}

		strippedExpectedGeneratedHtml := strings.Join(strings.Fields(string(b.String())),"")

		w := httptest.NewRecorder()
		defer os.Chdir(ChangeDirForTest())
		if err := generateHTML(w, threads, "layout.html", "public.navbar.html","index.html"); err != nil{
			t.Errorf("valid inputs were not expecting error, got %s", err)
		}

		BodyContent := w.Body.Bytes()
		strippedBodyContent := strings.Join(strings.Fields(string(BodyContent)),"")

		if strippedExpectedGeneratedHtml != strippedBodyContent{
			t.Errorf("expected body content and got did not match.")
			t.Errorf("expected:%s", string(strippedExpectedGeneratedHtml))
			t.Errorf("got:%s", string(strippedBodyContent))
		}
	})

}