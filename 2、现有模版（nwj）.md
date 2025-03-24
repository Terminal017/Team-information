
### 前言：不会写的selector或者筛选添加可以问AI

使用下面模版时，请不要将注释拷进去。
```go

//根据sitemap采集
engine.OnXML("//loc", func(element *colly.XMLElement, ctx *crawlers.Context) {
	engine.Visit(element.Text, crawlers.News)
})
    
//根据特定字符筛选
engine.OnXML("//loc", func(element *colly.XMLElement, ctx *crawlers.Context) {
	if strings.Contains(element.Text, ".xml") {
		engine.Visit(element.Text, crawlers.Index)
	}
	if strings.Contains(element.Text, ".ece") {
		engine.Visit(element.Text, crawlers.News)
	}
})


//获取图片
engine.OnHTML("selector", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	ctx.Image = []string{element.Attr("src")}
})

//获取图片(并补全URL)
engine.OnHTML("selector > img", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	imageUrl := element.Request.AbsoluteURL(element.Attr("src"))
	ctx.Image = []string{imageUrl}
})


//获取副标题
engine.OnHTML("section", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	ctx.SubTitle += element.Text
})

//获取Description
engine.OnHTML("section", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	ctx.Description += element.Text
})

//获取新闻作者
engine.OnHTML("selector", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	ctx.Authors = append(ctx.Authors, strings.TrimSpace(element.Text))
})


//获取新闻时间(建议如果有时间且自动采采不到时，采一下时间)
engine.OnHTML("selector", func(element *colly.HTMLElement, ctx *crawlers.Context) {
    ctx.PublicationTime = strings.TrimSpace(element.Text)
})


//补全URL的写法
engine.OnHTML("selector", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	url, err := element.Request.URL.Parse(element.Attr("href"))
	if err != nil {
		crawlers.Sugar.Error(err.Error())
		return
	}
	engine.Visit(url.String(), crawlers.Index)
})

//过滤标签内不需要的标签(只能过滤标签的直接内容)
engine.OnHTML("selector > p",             //只对p的子元素生效
	func(element *colly.HTMLElement, ctx *crawlers.Context) {
		// 过滤noscript标签内容,避免可能的垃圾信息
		directText := element.DOM.Contents().Not("noscript").Text()
		ctx.Content += directText
	})

//通过清除的方式，深度过滤标签内不需要的标签(包括子元素的子元素...)
	engine.OnHTML("div.entry-content > p", func(element *colly.HTMLElement, ctx *crawlers.Context) {
	// 移除 p 标签中的所有 noscript 标签
	element.DOM.Find("noscript").Remove()
	directText := element.DOM.Text()
	ctx.Content += directText
})


//采集PDF
	engine.OnHTML("selector", func(element *colly.HTMLElement, ctx *crawlers.Context) {
		fileURL := element.Attr("href")
		if strings.Contains(fileURL, ".pdf") {        //更改这里和下面以修改采集的文件类型
			url, err := element.Request.URL.Parse(element.Attr("href"))
			if err != nil {
				crawlers.Sugar.Error(err.Error())
				return
			}
			ctx.File = append(ctx.File, url.String())      //可以改成video以采集视频链接
			ctx.PageType = crawlers.Report
		}
	})
```