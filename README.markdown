HtmlSuds
========
*Sqeaky-clean HTML markup*

**Disclaimer:** HtmlSuds is in the testing stages. The syntax may change, and it might not work at all for you. If it doesn't *compile*, that's probably a bug, and you should file it as such. However, if it just doesn't *work*, or doesn't quite work *properly*, then that is to be expected at this stage.

**Discussions:** If you have any suggestions for changing the syntax or underlying features of HtmlSuds, file an issue and we can have a discussion there. It's very much up for debate at the moment, and I welcome suggestions and improvements.

HtmlSuds is an experiment in designing a flexible, lightweight, and easily extensible html preprocessing library. It is designed to be 99% backwards-compatible with regular html, so you can mix and match as required. For example, the following is completely valid:

	<html>
		<body>
			@markdown
	Here, I can write the body of my site in markdown! This is really helpful when writing everything from help pages to simple descriptions for form elements. **Awesome**.
			@!markdown
		</body>
	</html>

HtmlSuds was designed to solve problems with normal HTML that make writing DRY code difficult. Because there are a lot of these, I've broken them into a few different sections.

Simpler Script and StyleSheet Includes
--------------------------------------
I don't know about you, but writing

	<script type="text/javascript" src="/static/js/foo.js"></script>

seems like an awful lot of cruft to me. Let's see what this looks like in HtmlSuds:

	@script src="/static/js/foo.js" !

If all of your assets are coming from the same folder or server, you could easily customize HtmlSuds to add that part automatically:

	@script src="foo.js" !

Much better, right? Styles are similarilly improved, going from 

	<link rel="stylesheet" type="text/css" href="/static/css/foo.css" />

down to (with the customization explained above) the much simpler

	@link href="foo.css" !

However, the best part about this feature is not the simpler syntax and the DRYer code, but the flexibility. Now that the preprocessor knows what exactly it is looking at, it can make optimizations easily. For example, if foo.js was a 1K javascript file, it might choose to inline it in order to increase performance. This allows you to keep your script files well structured, without worrying about performance issues. If you encounter a performance issue later in your development cycle, you can simply tweak the preprocessor to work better in your particular situation.

**Documentation (and project) in progress. More to come.**