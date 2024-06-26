// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF60e8083DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(in *jlexer.Lexer, out *PowProof) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = uint32(in.Uint32())
		case "uuid":
			out.UUID = string(in.String())
		case "access_token_uuid":
			out.AccessTokenUUID = string(in.String())
		case "message_check_nonce":
			out.MessageCheckNonce = int64(in.Int64())
		case "message_hash":
			out.MessageHash = string(in.String())
		case "message_data":
			if in.IsNull() {
				in.Skip()
				out.MessageData = nil
			} else {
				out.MessageData = in.Bytes()
			}
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updated_at":
			if in.IsNull() {
				in.Skip()
				out.UpdatedAt = nil
			} else {
				if out.UpdatedAt == nil {
					out.UpdatedAt = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.UpdatedAt).UnmarshalJSON(data))
				}
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF60e8083EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(out *jwriter.Writer, in PowProof) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	{
		const prefix string = ",\"uuid\":"
		out.RawString(prefix)
		out.String(string(in.UUID))
	}
	{
		const prefix string = ",\"access_token_uuid\":"
		out.RawString(prefix)
		out.String(string(in.AccessTokenUUID))
	}
	{
		const prefix string = ",\"message_check_nonce\":"
		out.RawString(prefix)
		out.Int64(int64(in.MessageCheckNonce))
	}
	{
		const prefix string = ",\"message_hash\":"
		out.RawString(prefix)
		out.String(string(in.MessageHash))
	}
	{
		const prefix string = ",\"message_data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.MessageData)
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	if in.UpdatedAt != nil {
		const prefix string = ",\"updated_at\":"
		out.RawString(prefix)
		out.Raw((*in.UpdatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PowProof) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF60e8083EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PowProof) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF60e8083EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PowProof) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF60e8083DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PowProof) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF60e8083DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(l, v)
}
