package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// pop3 package

const contentTypeMultipartMixed = "multipart/mixed"
const contentTypeMultipartAlternative = "multipart/alternative"
const contentTypeMultipartRelated = "multipart/related"
const contentTypeTextHtml = "text/html"
const contentTypeTextPlain = "text/plain"

// Client for POP3.
type Client struct {
	Text      *Connection
	conn      net.Conn
	timeout   time.Duration
	useTLS    bool
	tlsConfig *tls.Config
}

// MessageInfo represents the message attributes returned by a LIST command.
type MessageInfo struct {
	Seq  uint32 // Message sequence number
	Size uint32 // Message size in bytes
	UID  string // Message UID
}

type option func(*Client) option

// Noop is a configuration function that does nothing.
func Noop() option {
	return func(c *Client) option {
		return Noop()
	}
}

// UseTLS is a configuration function whose result is passed as a parameter in
// the Dial function. It configures the client to use TLS.
func UseTLS(config *tls.Config) option {
	return func(c *Client) option {
		c.useTLS = true
		c.tlsConfig = config
		return Noop()
	}
}

// UseTimeout is a configuration function whose result is passed as a parameter in
// the Dial function. It configures the client to use timeouts for each POP3 command.
func UseTimeout(timeout time.Duration) option {
	return func(c *Client) option {
		previous := c.timeout
		c.UseTimeouts(timeout)
		return UseTimeout(previous)
	}
}

const (
	protocol      = "tcp"
	lineSeparator = "\n"
)

// Dial connects to the given address and returns a client holding a tcp connection.
// To pass configuration to the Dial function use the methods UseTLS or UseTimeout.
// E.g. c, err = pop3.Dial(address, pop3.UseTLS(tlsConfig), pop3.UseTimeout(timeout))
func Dial(addr string, options ...option) (*Client, error) {
	client := &Client{}
	for _, option := range options {
		option(client)
	}
	var (
		conn net.Conn
		err  error
	)
	if !client.useTLS {
		if client.timeout > time.Duration(0) {
			conn, err = net.DialTimeout(protocol, addr, client.timeout)
		} else {
			conn, err = net.Dial(protocol, addr)
		}
	} else {
		host, _, _ := net.SplitHostPort(addr)
		if client.timeout > time.Duration(0) {
			d := net.Dialer{Timeout: client.timeout}
			conn, err = tls.DialWithDialer(&d, protocol, addr, setServerName(client.tlsConfig, host))
		} else {
			conn, err = tls.Dial(protocol, addr, setServerName(client.tlsConfig, host))

		}
	}
	if err != nil {
		return nil, err
	}
	client.conn = conn
	err = client.initialize()
	if err != nil {
		return nil, err
	}
	return client, nil

}

// NewClient initializeds a client.
// To pass configuration to the NewClient function use the methods UseTLS or UseTimeout.
// E.g. c, err = pop3.Dial(address, pop3.UseTLS(tlsConfig), pop3.UseTimeout(timeout))
func NewClient(conn net.Conn, options ...option) (*Client, error) {
	client := &Client{conn: conn}

	for _, option := range options {
		option(client)
	}

	err := client.initialize()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) initialize() (err error) {
	text := NewConnection(client.conn)
	client.Text = text
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.ReadResponse()
	return
}

// UseTimeouts adds a timeout to the client. Timeouts are applied on every
// POP3 command.
func (client *Client) UseTimeouts(timeout time.Duration) {
	client.timeout = timeout
}

// User issues the POP3 User command.
func (client *Client) User(user string) (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("USER %s", user)
	return
}

// Pass sends the given password to the server. The password is sent
// unencrypted unless the connection is already secured by TLS (via DialTLS or
// some other mechanism).
func (client *Client) Pass(password string) (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("PASS %s", password)
	return
}

// Auth sends the given username and password to the server.
func (client *Client) Auth(username, password string) (err error) {
	err = client.User(username)
	if err != nil {
		return
	}
	err = client.Pass(password)
	return
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// Information provided besides the number of messages and the size of the
// maildrop is ignored. In the event of an error, all returned numeric values
// will be 0.
func (client *Client) Stat() (count, size uint32, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	l, err := client.Text.Cmd("STAT")
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(l)
	count, err = stringToUint32(parts[0])
	if err != nil {
		return 0, 0, errors.New("Invalid server response")
	}
	size, err = stringToUint32(parts[1])
	if err != nil {
		return 0, 0, errors.New("Invalid server response")
	}
	return
}

// List returns the size of the message referenced by the sequence number,
// if it exists. If the message does not exist, or another error is encountered,
// the returned size will be 0.
func (client *Client) List(msgSeqNum uint32) (size uint32, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	l, err := client.Text.Cmd("LIST %d", msgSeqNum)
	if err != nil {
		return 0, err
	}
	size, err = stringToUint32(strings.Fields(l)[1])
	if err != nil {
		return 0, errors.New("Invalid server response")
	}
	return size, nil
}

// ListAll returns a list of MessageInfo for all messages, containing their
// sequence number and size.
func (client *Client) ListAll() (msgInfos []*MessageInfo, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("LIST")
	if err != nil {
		return
	}
	lines, err := client.Text.ReadMultiLines()
	if err != nil {

		return
	}
	msgInfos = make([]*MessageInfo, len(lines))
	for i, line := range lines {
		var seq, size uint32
		fields := strings.Fields(line)
		seq, err = stringToUint32(fields[0])
		if err != nil {
			return
		}
		size, err = stringToUint32(fields[1])
		if err != nil {
			return
		}

		msgInfos[i] = &MessageInfo{
			Seq:  seq,
			Size: size,
		}

	}
	return
}

// Retr downloads and returns the given message. The lines are separated by LF,
// whatever the server sent.
func (client *Client) Retr(msg uint32) (text string, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("RETR %d", msg)
	if err != nil {
		return "", err
	}
	lines, err := client.Text.ReadMultiLines()
	text = strings.Join(lines, lineSeparator)
	return
}

func (client *Client) Top(msg uint32, msgNumStr uint32) (text string, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("TOP %d %d", msg, msgNumStr)
	if err != nil {
		return "", err
	}
	lines, err := client.Text.ReadMultiLines()
	text = strings.Join(lines, lineSeparator)
	return
}

// Dele marks the given message as deleted.
func (client *Client) Dele(msg uint32) (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("DELE %d", msg)
	return
}

// Noop does nothing, but will prolong the end of the connection if the server
// has a timeout set.
func (client *Client) Noop() (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("NOOP")
	return
}

// Rset unmarks any messages marked for deletion previously in this session.
func (client *Client) Rset() (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("RSET")
	return
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (client *Client) Quit() (err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("QUIT")
	if err != nil {
		return err
	}
	client.Text.Close()
	return
}

// UIDl retrieves the unique ID of the message referenced by the sequence number.
func (client *Client) UIDl(msgSeqNum uint32) (uid string, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	line, err := client.Text.Cmd("UIDL %d", msgSeqNum)
	if err != nil {
		return "", err
	}
	uid = strings.Fields(line)[1]
	return
}

// UIDlAll retrieves the unique IDs and sequence number for all messages.
func (client *Client) UIDlAll() (msgInfos []*MessageInfo, err error) {
	client.setDeadline()
	defer client.resetDeadline()
	_, err = client.Text.Cmd("UIDL")
	if err != nil {
		return
	}
	lines, err := client.Text.ReadMultiLines()
	if err != nil {
		return
	}
	msgInfos = make([]*MessageInfo, len(lines))
	for i, line := range lines {
		var seq uint32
		var uid string
		fields := strings.Fields(line)
		seq, err = stringToUint32(fields[0])
		if err != nil {
			return
		}
		uid = fields[1]
		msgInfos[i] = &MessageInfo{
			Seq: seq,
			UID: uid,
		}
	}
	return
}

func (client *Client) setDeadline() {
	if client.timeout > time.Duration(0) {
		client.conn.SetDeadline(time.Now().Add(client.timeout))
	}
}

func (client *Client) resetDeadline() {
	if client.timeout > time.Duration(0) {
		client.conn.SetDeadline(time.Time{})
	}
}

func stringToUint32(intString string) (uint32, error) {
	val, err := strconv.Atoi(intString)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

// setServerName returns a new TLS configuration with ServerName set to host if
// the original configuration was nil or config.ServerName was empty.
// Copied from go-imap: code.google.com/p/go-imap/go1/imap
func setServerName(config *tls.Config, host string) *tls.Config {
	if config == nil {
		config = &tls.Config{ServerName: host}
	} else if config.ServerName == "" {
		c := *config
		c.ServerName = host
		config = &c
	}
	return config
}


// Parse an email message read from io.Reader into parsemail.Email struct
func ShortParse(r io.Reader) (email Email, err error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return
	}

	email, err = createShortEmailFromHeader(msg.Header)
	if err != nil {
		return
	}

	return
}

// Parse an email message read from io.Reader into parsemail.Email struct
func Parse(r io.Reader) (email Email, err error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return
	}

	email, err = createEmailFromHeader(msg.Header)
	if err != nil {
		return
	}

	contentType, params, err := parseContentType(msg.Header.Get("Content-Type"))
	if err != nil {
		return
	}

	switch contentType {
	case contentTypeMultipartMixed:
		email.TextBody, email.HTMLBody, email.Attachments, email.EmbeddedFiles, err = parseMultipartMixed(msg.Body, params["boundary"])
	case contentTypeMultipartAlternative:
		email.TextBody, email.HTMLBody, email.EmbeddedFiles, err = parseMultipartAlternative(msg.Body, params["boundary"])
	case contentTypeTextPlain:
		message, _ := ioutil.ReadAll(msg.Body)
		email.TextBody = strings.TrimSuffix(string(message[:]), "\n")
	case contentTypeTextHtml:
		message, _ := ioutil.ReadAll(msg.Body)
		email.HTMLBody = strings.TrimSuffix(string(message[:]), "\n")
	default:
		err = fmt.Errorf("Unknown top level mime type: %s", contentType)
	}

	return
}

// Парсинг письма из команды TOP
func createShortEmailFromHeader(header mail.Header) (email Email, err error) {
	hp := headerParser{header: &header}

	email.Subject = decodeMimeSentence(header.Get("Subject"))
	email.From = hp.parseAddressList(header.Get("From"))
	email.Sender = hp.parseAddress(header.Get("Sender"))
	email.To = hp.parseAddressList(header.Get("To"))
	email.Date = hp.parseTime(header.Get("Date"))
	email.MessageID = hp.parseMessageId(header.Get("Message-ID"))

	if hp.err != nil {
		err = hp.err
		return
	}

	//decode whole header for easier access to extra fields
	//todo: should we decode? aren't only standard fields mime encoded?
	email.Header, err = decodeHeaderMime(header)
	if err != nil {
		return
	}

	return
}

func createEmailFromHeader(header mail.Header) (email Email, err error) {
	hp := headerParser{header: &header}

	email.Subject = decodeMimeSentence(header.Get("Subject"))
	email.From = hp.parseAddressList(header.Get("From"))
	email.Sender = hp.parseAddress(header.Get("Sender"))
	email.ReplyTo = hp.parseAddressList(header.Get("Reply-To"))
	email.To = hp.parseAddressList(header.Get("To"))
	email.Cc = hp.parseAddressList(header.Get("Cc"))
	email.Bcc = hp.parseAddressList(header.Get("Bcc"))
	email.Date = hp.parseTime(header.Get("Date"))
	email.ResentFrom = hp.parseAddressList(header.Get("Resent-From"))
	email.ResentSender = hp.parseAddress(header.Get("Resent-Sender"))
	email.ResentTo = hp.parseAddressList(header.Get("Resent-To"))
	email.ResentCc = hp.parseAddressList(header.Get("Resent-Cc"))
	email.ResentBcc = hp.parseAddressList(header.Get("Resent-Bcc"))
	email.ResentMessageID = hp.parseMessageId(header.Get("Resent-Message-ID"))
	email.MessageID = hp.parseMessageId(header.Get("Message-ID"))
	email.InReplyTo = hp.parseMessageIdList(header.Get("In-Reply-To"))
	email.References = hp.parseMessageIdList(header.Get("References"))
	email.ResentDate = hp.parseTime(header.Get("Resent-Date"))

	if hp.err != nil {
		err = hp.err
		return
	}

	//decode whole header for easier access to extra fields
	//todo: should we decode? aren't only standard fields mime encoded?
	email.Header, err = decodeHeaderMime(header)
	if err != nil {
		return
	}

	return
}

func parseContentType(contentTypeHeader string) (contentType string, params map[string]string, err error) {
	if contentTypeHeader == "" {
		contentType = contentTypeTextPlain
		return
	}

	return mime.ParseMediaType(contentTypeHeader)
}

func parseMultipartRelated(msg io.Reader, boundary string) (textBody, htmlBody string, embeddedFiles []EmbeddedFile, err error) {
	pmr := multipart.NewReader(msg, boundary)
	for {
		part, err := pmr.NextPart()

		if err == io.EOF {
			break
		} else if err != nil {
			return textBody, htmlBody, embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return textBody, htmlBody, embeddedFiles, err
		}

		switch contentType {
		case contentTypeTextPlain:
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			textBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		case contentTypeTextHtml:
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			htmlBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		case contentTypeMultipartAlternative:
			tb, hb, ef, err := parseMultipartAlternative(part, params["boundary"])
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			htmlBody += hb
			textBody += tb
			embeddedFiles = append(embeddedFiles, ef...)
		default:
			if isEmbeddedFile(part) {
				ef, err := decodeEmbeddedFile(part)
				if err != nil {
					return textBody, htmlBody, embeddedFiles, err
				}

				embeddedFiles = append(embeddedFiles, ef)
			} else {
				return textBody, htmlBody, embeddedFiles, fmt.Errorf("Can't process multipart/related inner mime type: %s", contentType)
			}
		}
	}

	return textBody, htmlBody, embeddedFiles, err
}

func parseMultipartAlternative(msg io.Reader, boundary string) (textBody, htmlBody string, embeddedFiles []EmbeddedFile, err error) {
	pmr := multipart.NewReader(msg, boundary)
	for {
		part, err := pmr.NextPart()

		if err == io.EOF {
			break
		} else if err != nil {
			return textBody, htmlBody, embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return textBody, htmlBody, embeddedFiles, err
		}

		switch contentType {
		case contentTypeTextPlain:
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			textBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		case contentTypeTextHtml:
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			htmlBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		case contentTypeMultipartRelated:
			tb, hb, ef, err := parseMultipartRelated(part, params["boundary"])
			if err != nil {
				return textBody, htmlBody, embeddedFiles, err
			}

			htmlBody += hb
			textBody += tb
			embeddedFiles = append(embeddedFiles, ef...)
		default:
			if isEmbeddedFile(part) {
				ef, err := decodeEmbeddedFile(part)
				if err != nil {
					return textBody, htmlBody, embeddedFiles, err
				}

				embeddedFiles = append(embeddedFiles, ef)
			} else {
				return textBody, htmlBody, embeddedFiles, fmt.Errorf("Can't process multipart/alternative inner mime type: %s", contentType)
			}
		}
	}

	return textBody, htmlBody, embeddedFiles, err
}

func parseMultipartMixed(msg io.Reader, boundary string) (textBody, htmlBody string, attachments []Attachment, embeddedFiles []EmbeddedFile, err error) {
	mr := multipart.NewReader(msg, boundary)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return textBody, htmlBody, attachments, embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return textBody, htmlBody, attachments, embeddedFiles, err
		}

		if contentType == contentTypeMultipartAlternative {
			textBody, htmlBody, embeddedFiles, err = parseMultipartAlternative(part, params["boundary"])
			if err != nil {
				return textBody, htmlBody, attachments, embeddedFiles, err
			}
		} else if contentType == contentTypeTextPlain {
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, attachments, embeddedFiles, err
			}

			textBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		} else if contentType == contentTypeTextHtml {
			ppContent, err := ioutil.ReadAll(part)
			if err != nil {
				return textBody, htmlBody, attachments, embeddedFiles, err
			}

			textBody += strings.TrimSuffix(string(ppContent[:]), "\n")
		} else if contentType == contentTypeMultipartRelated {
			textBody, htmlBody, embeddedFiles, err = parseMultipartRelated(part, params["boundary"])
			if err != nil {
				return textBody, htmlBody, attachments, embeddedFiles, err
			}
		} else if isAttachment(part) {
			at, err := decodeAttachment(part)
			if err != nil {
				return textBody, htmlBody, attachments, embeddedFiles, err
			}

			attachments = append(attachments, at)
		} else {
			return textBody, htmlBody, attachments, embeddedFiles, fmt.Errorf("Unknown multipart/mixed nested mime type: %s", contentType)
		}
	}

	return textBody, htmlBody, attachments, embeddedFiles, err
}

func decodeMimeSentence(s string) string {
	result := []string{}
	ss := strings.Split(s, " ")

	for _, word := range ss {
		dec := new(mime.WordDecoder)
		w, err := dec.Decode(word)
		if err != nil {
			if len(result) == 0 {
				w = word
			} else {
				w = " " + word
			}
		}

		result = append(result, w)
	}

	return strings.Join(result, "")
}

func decodeHeaderMime(header mail.Header) (mail.Header, error) {
	parsedHeader := map[string][]string{}

	for headerName, headerData := range header {

		parsedHeaderData := []string{}
		for _, headerValue := range headerData {
			parsedHeaderData = append(parsedHeaderData, decodeMimeSentence(headerValue))
		}

		parsedHeader[headerName] = parsedHeaderData
	}

	return mail.Header(parsedHeader), nil
}

func decodePartData(part *multipart.Part) (io.Reader, error) {
	encoding := part.Header.Get("Content-Transfer-Encoding")

	if strings.EqualFold(encoding, "base64") {
		dr := base64.NewDecoder(base64.StdEncoding, part)
		dd, err := ioutil.ReadAll(dr)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(dd), nil
	}

	return nil, fmt.Errorf("Unknown encoding: %s", encoding)
}

func isEmbeddedFile(part *multipart.Part) bool {
	return part.Header.Get("Content-Transfer-Encoding") != ""
}

func decodeEmbeddedFile(part *multipart.Part) (ef EmbeddedFile, err error) {
	cid := decodeMimeSentence(part.Header.Get("Content-Id"))
	decoded, err := decodePartData(part)
	if err != nil {
		return
	}

	ef.CID = strings.Trim(cid, "<>")
	ef.Data = decoded
	ef.ContentType = part.Header.Get("Content-Type")

	return
}

func isAttachment(part *multipart.Part) bool {
	return part.FileName() != ""
}

func decodeAttachment(part *multipart.Part) (at Attachment, err error) {
	filename := decodeMimeSentence(part.FileName())
	decoded, err := decodePartData(part)
	if err != nil {
		return
	}

	at.Filename = filename
	at.Data = decoded
	at.ContentType = strings.Split(part.Header.Get("Content-Type"), ";")[0]

	return
}

type headerParser struct {
	header *mail.Header
	err    error
}

func (hp headerParser) parseAddress(s string) (ma *mail.Address) {
	if hp.err != nil {
		return nil
	}

	if strings.Trim(s, " \n") != "" {
		ma, hp.err = mail.ParseAddress(s)

		return ma
	}

	return nil
}

func (hp headerParser) parseAddressList(s string) (ma []*mail.Address) {
	if hp.err != nil {
		return
	}

	if strings.Trim(s, " \n") != "" {
		ma, hp.err = mail.ParseAddressList(s)
		return
	}

	return
}

func (hp headerParser) parseTime(s string) (t time.Time) {
	if hp.err != nil || s == "" {
		return
	}
	regexpString := " \\(.*\\)$"
	re := regexp.MustCompile(regexpString)
	findedStr := re.FindString(s)
	newString := ""
	if len(findedStr)>0 {
		newString = strings.Replace(s,findedStr,"",-1)
	} else {
		newString = s
	}
	t, hp.err = time.Parse(time.RFC1123Z, newString)
	if hp.err == nil {
		return t
	}

	t, hp.err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", newString)

	return
}

func (hp headerParser) parseMessageId(s string) string {
	if hp.err != nil {
		return ""
	}

	return strings.Trim(s, "<> ")
}

func (hp headerParser) parseMessageIdList(s string) (result []string) {
	if hp.err != nil {
		return
	}

	for _, p := range strings.Split(s, " ") {
		if strings.Trim(p, " \n") != "" {
			result = append(result, hp.parseMessageId(p))
		}
	}

	return
}

// Attachment with filename, content type and data (as a io.Reader)
type Attachment struct {
	Filename    string
	ContentType string
	Data        io.Reader
}

// EmbeddedFile with content id, content type and data (as a io.Reader)
type EmbeddedFile struct {
	CID         string
	ContentType string
	Data        io.Reader
}

// Email with fields for all the headers defined in RFC5322 with it's attachments and
type Email struct {
	Header mail.Header

	Subject    string
	Sender     *mail.Address
	From       []*mail.Address
	ReplyTo    []*mail.Address
	To         []*mail.Address
	Cc         []*mail.Address
	Bcc        []*mail.Address
	Date       time.Time
	MessageID  string
	InReplyTo  []string
	References []string

	ResentFrom      []*mail.Address
	ResentSender    *mail.Address
	ResentTo        []*mail.Address
	ResentDate      time.Time
	ResentCc        []*mail.Address
	ResentBcc       []*mail.Address
	ResentMessageID string

	HTMLBody string
	TextBody string

	Attachments   []Attachment
	EmbeddedFiles []EmbeddedFile
}

// Connection stores a Reader and a Writer.
type Connection struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	conn   io.ReadWriteCloser
}

var crlf = []byte{'\r', '\n'}
var okResponse = "+OK"
var endResponse = "."

// NewConnection initializes a connection.
func NewConnection(conn io.ReadWriteCloser) *Connection {
	return &Connection{
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
		conn:   conn,
	}
}

// Close closes a connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Cmd sends the given command on the connection.
func (c *Connection) Cmd(format string, args ...interface{}) (result string, err error) {
	c.SendCMD(format, args...)
	return c.ReadResponse()
}

// SendCMD writes the command on the writer and flushes the writer afterwards.
func (c *Connection) SendCMD(format string, args ...interface{}) {
	fmt.Fprintf(c.Writer, format, args...)
	c.Writer.Write(crlf)
	c.Writer.Flush()

	return
}

// ReadResponse reads the response from the server and parses it.
// It checks whether the response is OK and returns the result omitting the OK+ prefix.
func (c *Connection) ReadResponse() (result string, err error) {
	result = ""

	response, _, err := c.Reader.ReadLine()
	if err != nil {
		return
	}

	line := string(response)
	if line[0:3] != okResponse {

		err = errors.New(line[5:])
	}

	if len(line) >= 4 {
		result = line[4:]
	}

	return
}

// ReadMultiLines reads a response with multiple lines.
func (c *Connection) ReadMultiLines() (lines []string, err error) {
	lines = make([]string, 0)
	var bytes []byte

	for {
		bytes, _, err = c.Reader.ReadLine()
		line := string(bytes)

		if err != nil || line == endResponse {
			return
		}

		if len(line) > 0 && string(line[0]) == "." {
			line = line[1:]
		}

		lines = append(lines, line)
	}

	return
}
