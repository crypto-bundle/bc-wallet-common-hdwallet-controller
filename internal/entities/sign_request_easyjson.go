// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

import (
	json "encoding/json"
	types "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pq "github.com/lib/pq"
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

func easyjsonE949d281DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(in *jlexer.Lexer, out *SignRequest) {
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
		case "wallet_uuid":
			out.WalletUUID = string(in.String())
		case "session_uuid":
			out.SessionUUID = string(in.String())
		case "purpose_uuid":
			out.PurposeUUID = string(in.String())
		case "status":
			out.Status = types.SignRequestStatus(in.Uint8())
		case "derivation_path":
			if in.IsNull() {
				in.Skip()
				out.DerivationPath = nil
			} else {
				in.Delim('[')
				if out.DerivationPath == nil {
					if !in.IsDelim(']') {
						out.DerivationPath = make(pq.Int32Array, 0, 16)
					} else {
						out.DerivationPath = pq.Int32Array{}
					}
				} else {
					out.DerivationPath = (out.DerivationPath)[:0]
				}
				for !in.IsDelim(']') {
					var v1 int32
					v1 = int32(in.Int32())
					out.DerivationPath = append(out.DerivationPath, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonE949d281EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(out *jwriter.Writer, in SignRequest) {
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
		const prefix string = ",\"wallet_uuid\":"
		out.RawString(prefix)
		out.String(string(in.WalletUUID))
	}
	{
		const prefix string = ",\"session_uuid\":"
		out.RawString(prefix)
		out.String(string(in.SessionUUID))
	}
	{
		const prefix string = ",\"purpose_uuid\":"
		out.RawString(prefix)
		out.String(string(in.PurposeUUID))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Uint8(uint8(in.Status))
	}
	if len(in.DerivationPath) != 0 {
		const prefix string = ",\"derivation_path\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v2, v3 := range in.DerivationPath {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.Int32(int32(v3))
			}
			out.RawByte(']')
		}
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
func (v SignRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE949d281EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SignRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE949d281EncodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SignRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE949d281DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SignRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE949d281DecodeGithubComCryptoBundleBcWalletCommonHdwalletControllerInternalEntities(l, v)
}
