// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

import (
	json "encoding/json"
	types "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
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

func easyjson5eca2d97DecodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(in *jlexer.Lexer, out *MnemonicWalletSession) {
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
			out.AccessTokenUUID = uint32(in.Uint32())
		case "mnemonic_wallet_uuid":
			out.MnemonicWalletUUID = string(in.String())
		case "status":
			out.Status = types.MnemonicWalletSessionStatus(in.Uint8())
		case "expired_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ExpiredAt).UnmarshalJSON(data))
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
func easyjson5eca2d97EncodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(out *jwriter.Writer, in MnemonicWalletSession) {
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
		out.Uint32(uint32(in.AccessTokenUUID))
	}
	{
		const prefix string = ",\"mnemonic_wallet_uuid\":"
		out.RawString(prefix)
		out.String(string(in.MnemonicWalletUUID))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Uint8(uint8(in.Status))
	}
	{
		const prefix string = ",\"expired_at\":"
		out.RawString(prefix)
		out.Raw((in.ExpiredAt).MarshalJSON())
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"updated_at\":"
		out.RawString(prefix)
		if in.UpdatedAt == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.UpdatedAt).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MnemonicWalletSession) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5eca2d97EncodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MnemonicWalletSession) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5eca2d97EncodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MnemonicWalletSession) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5eca2d97DecodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MnemonicWalletSession) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5eca2d97DecodeGithubComCryptoBundleBcWalletCommonHdwalletManagerInternalEntities(l, v)
}
