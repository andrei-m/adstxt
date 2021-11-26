# adstxt
[![GoDoc](https://godoc.org/github.com/andrei-m/adstxt?status.svg)](https://godoc.org/github.com/andrei-m/adstxt)

ads.txt resolver for Go. This implementation is based on v1.03 of the [ads.txt spec](https://iabtechlab.com/wp-content/uploads/2020/12/ads-txt-v1.0.3_Draft_in_Public_Comment_IABTechLab_2020-12.pdf) as published by the IAB Tech Lab.

## Examples

Resolve ads.txt from a domain using the default HTTP client:

```
adstxt, err := adstxt.DefaultResolve("www.example.com")
```

Customize the HTTP client. Consider using this package's CheckRedirect function, because it implements the IAB-specified external redirect rules.

```
client := &http.Client{
	Timeout: 1 * time.Minute,
	CheckRedirect: adstxt.CheckRedirect,
}
adstxt, err := adstxt.Resolve(client, "www.example.com")
```

Parse ads.txt data without making an HTTP request

```
rawAdsTxt := strings.NewReader(`# comment
foo,bar,DIRECT,baz
three,four,RESELLER`)
adstxt, err := adstxt.Parse(rawAdsTxt)
```

See the GoDoc link for a description of the parsed ads.txt format.

