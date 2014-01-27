package Rserve

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func serialize(method, commandStr string) {
	if method == "" {
		if strings.HasSuffix(commandStr, ";") {
			method = "voidEval"
		} else {
			method = "eval"
		}
	}
	strlen := len([]byte(commandStr))
	if strlen%4 == 0 {
		strlen += 4 - (strlen % 4) // Not even sure why this is a thing
		// I've gathered it is to "Ensure it's a multiple of 4"
	}
	buf := make([]byte, 16+4+strlen)
	for i := 0; i < len(buf); i++ {
		buf[i] = 0x00
	}

	cmdcode := uint32(getcommandcode(method))

	//     buf.writeUInt32LE(cmdCode, 0); // Command code
	tbuf := new(bytes.Buffer)
	err := binary.Write(tbuf, binary.LittleEndian, cmdcode)
	if err != nil {
		panic("wat")
	}
	uintbuf := make([]byte, 4)
	n, err := tbuf.Read(uintbuf)
	if n != 4 || err != nil {
		panic("wat man")
	}
	for k, _ := range uintbuf {
		buf[k] = uintbuf[k]
	}
	//     buf.writeUInt32LE(4 + strlen, 4); // data length

}

// serialize = function(method, commandStr) {
//     if (_.isUndefined(commandStr)) {
//         commandStr = method;
//         // Use voidEval when no method specified, and string ends in semicolon
//         method = (commandStr.match(/;(\n|\r)*$/)) ? 'voidEval' : 'eval';
//     }
//     var strlen = Buffer.byteLength(commandStr);
//     if (strlen % 4) {
//         strlen += 4 - (strlen % 4); // Ensure it's a multiple of 4, possibly not needed
//     }
//     var buf = new Buffer(16 + 4 + strlen);
//     buf.fill(0x00);
//     var cmdCode = get_command_code(method);
//     buf.writeUInt32LE(cmdCode, 0); // Command code
//     buf.writeUInt32LE(4 + strlen, 4); // data length
//     // FIXME: The following line should actually be righting a 24bit integers in the 17-19 offset
//     // Instead it's writing a 16bit integer in the 17-18 offset
//     buf.writeUInt16LE(strlen, 17); // Length of the string
//     buf.writeUInt8(0x04, 16);  // Data is a string
//     if (strlen > 0) {
//         buf.write(commandStr, 20, 'utf8');
//     }
//     return buf;
// }
