const crypto = require('crypto');
const axios = require('axios');

const token = 'SwitchBotトークン';
const secret = 'SwitchBotシークレット';

const nowDate = Date.now();
const nonce = crypto.randomUUID();

const data = token + nowDate + nonce;
const signTerm = crypto
  .createHmac('sha256', secret)
  .update(Buffer.from(data, 'utf-8'))
  .digest();
const sign = signTerm.toString('base64');

const headers = {
  Authorization: token,
  sign,
  nonce,
  t: nowDate,
  'Content-Type': 'application/json'
};

const body = {
  command: 'turnOn', // スイッチをOFFにする場合は 'turnOff' にする
  parameter: 'default',
  commandType: 'command'
};

try {
  axios
    .post(
      'https://api.switch-bot.com/v1.1/devices/デバイスのMACアドレス/commands',
      body,
      { headers }
    )
    .then((b) => {
      console.log(b.data);
    });
} catch (e) {
  console.log(e);
}
