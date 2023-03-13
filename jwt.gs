const createJwt = ({ privateKey, expiresInMinutes, data = {} }) => {
    // Sign token using HMAC with SHA-256 algorithm
    const header = {
      alg: 'HS256',
      typ: 'JWT',
    };
  
    const now = Date.now();
    const expires = new Date(now);
    expires.setMinutes(expires.getMinutes() + expiresInMinutes);
  
    // iat = issued time, exp = expiration time
    const payload = {
      exp: Math.round(expires.getTime() / 1000),
      iat: Math.round(now / 1000),
    };
  
    // add user payload
    Object.keys(data).forEach(function (key) {
      payload[key] = data[key];
    });
  
    const base64Encode = (text, json = true) => {
      const data = json ? JSON.stringify(text) : text;
      return Utilities.base64EncodeWebSafe(data).replace(/=+$/, '');
    };
  
    const toSign = `${base64Encode(header)}.${base64Encode(payload)}`;
    const signatureBytes = Utilities.computeHmacSha256Signature(toSign, privateKey);
    const signature = base64Encode(signatureBytes, false);
    return `${toSign}.${signature}`;
};

const generateAccessToken = () => {
    const privateKey = 'Q1Sji37e2NXr9iauLNFbrTrFsui7D6/0Is9yw/O+';
    const accessToken = createJwt({
      privateKey,
      expiresInMinutes: 10,
      data: {
        sid: scriptProp.getProperty('key'),
      },
    });
    return accessToken
};
