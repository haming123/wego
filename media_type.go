package wego

import (
	"strings"
)

//Content-Type：application/x-www-form-urlencoded; charset:utf-8;
//Content-Type：multipart/form-data; boundary=----WebKitFormBoundary7TMYhSONfkAM2z3a
func parseMediaType(v string) (string, string) {
	rest := v
	mediatype := ""
	i := strings.Index(v, ";")
	if i == -1 {
		mediatype = rest
		rest = ""
	} else {
		mediatype = rest[0:i]
		rest = rest[i+1:]
	}
	mediatype = trimSpace(strings.ToLower(mediatype))

	for len(rest) > 0 {
		i = strings.Index(rest, ";")
		part := ""
		if i == -1 {
			part = rest
			rest = ""
		} else{
			part = rest[0:i]
			rest = rest[i+1:]
		}

		i = strings.Index(part, "=")
		if i == -1 {
			continue
		}
		key := part[0:i]
		val := part[i+1:]

		key = trimSpace(strings.ToLower(key))
		if key == "boundary" {
			return mediatype, val
		}
	}

	return mediatype, ""
}

/*
func (this *FormParam)parsePostMultipart(bound string, maxMemory int64) error {
	// Reserve an additional 10 MB for non-file parts.
	maxValueBytes := maxMemory + int64(10<<20)
	if maxValueBytes <= 0 {
		if maxMemory < 0 {
			maxValueBytes = 0
		} else {
			maxValueBytes = math.MaxInt64
		}
	}

	r := this.ctx.Input.Request
	mr := multipart.NewReader(r.Body, bound)
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		name := p.FormName()
		if name == "" {
			continue
		}

		filename := p.FileName()
		if filename != "" {
			continue
		}

		var b bytes.Buffer
		n, err := io.CopyN(&b, p, maxValueBytes+1)
		if err != nil && err != io.EOF {
			return err
		}
		maxValueBytes -= n
		if maxValueBytes < 0 {
			return ErrMessageTooLarge
		}
		this.param = append(this.param, ParamItem{name, b.String()})
	}

	return nil
}*/
