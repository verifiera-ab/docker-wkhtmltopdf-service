# wkhtmltopdf as a web service

A dockerized webservice written in [Go](https://golang.org/) that uses [wkhtmltopdf](http://wkhtmltopdf.org/) to convert HTML into documents (images or pdf files).

## Usage

The service listens on port 3000 for POST requests on the root path (`/`). Any other method returns a `405 not allowed` status. Any other path returns a `404 not found` status.

The body should contain a JSON-encoded object containing the following parameters:

**url**: The URL of the page to convert.

**output**: The type of document to generate, can be either `jpg`, `png` or `pdf`. Defauts to `pdf` if not specified. Depending on the output type the appropriate binary is called.

**options**: A list of key-value arguments that are passed on to the appropriate `wkhtmltopdf` binary. Boolean values are interpreted as flag arguments (e.g.: `--greyscale`).

**cookies**: A list of key-value arguments that are passed on to the appropriate `wkhtmltopdf` binary as separate `cookie` arguments.

** Example: ** posting the following JSON:

```
{
  "url": "http://www.google.com",
  "output": "jpg",
  "options": {
    "margin-bottom": "1cm",
    "orientation": "Landscape",
    "grayscale": true
  },
  "cookies": {
    "foo": "bar",
    "baz": "foo"
  }
}
```

would be equivalent to executing the following in the command-line:

```
/usr/local/bin/wkhtmltoimage --margin-bottom 1cm --orientation Landscape --grayscale --cookie foo bar --cookie baz foo --format jpg -q http://www.google.com -
```

The `-` at the end of the command-line is so that the document contents are redirected to stdout so we can in turn redirect it to the web server's response stream.

For more information on available options, please check [wkhtmltopdf usage documentation](http://wkhtmltopdf.org/usage/wkhtmltopdf.txt).


## Known limitations / TODOs

This service doesn't currently accept or convert raw HTML. We found that it's difficult to render raw HTML when you have relative links in the page. There are known workarounds that we haven't yet tried (see [this stack overflow question](http://stackoverflow.com/questions/21775572/wkhtmltopdf-relative-paths-in-html-with-redirected-in-out-streams-wont-work) for more information).


## License

This code is released under the [MIT License](http://opensource.org/licenses/MIT).
