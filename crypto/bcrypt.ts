import {enc, lib, SHA512} from "crypto-js";

export type WordArray = lib.WordArray;

//
export function BcryptSimple(data: WordArray, salt: WordArray, loops: number): WordArray {
    let hash = SHA512(data);
    // eslint-disable-next-line no-param-reassign
    salt = SHA512(salt);
    hash = SHA512(hash.concat(salt));

    let tempHash = hash.clone();
    let tempSalt = salt.clone();
    const hashs: WordArray = lib.WordArray.create();
    for (let i = 0; i < 10; i += 1) {
        // eslint-disable-next-line no-bitwise
        tempHash = SHA512(tempHash.concat(lib.WordArray.create([i << 24], 1)));

        // eslint-disable-next-line no-bitwise
        tempSalt = SHA512(tempSalt.concat(lib.WordArray.create([i << 24], 1)));

        tempHash = SHA512(tempHash.concat(tempSalt));
        hashs.concat(tempHash);
    }

    hash = SHA512(hashs);

    const idxHash = hash.sigBytes / 4 / 2;
    let hash1 = lib.WordArray.create(hash.words.slice(0, idxHash));
    let hash2 = lib.WordArray.create(hash.words.slice(idxHash));

    const idxSalt = salt.sigBytes / 4 / 2;
    const salt1 = lib.WordArray.create(salt.words.slice(0, idxSalt));
    const salt2 = lib.WordArray.create(salt.words.slice(idxSalt));

    hash1.clamp();
    hash2.clamp();

    for (let i = 0; i < loops; i += 1) {
        hash1 = SHA512(hash1.concat(salt1));
        hash2 = SHA512(hash2.concat(salt2));
    }

    hash = SHA512(hash1.concat(hash2));
    hash = SHA512(hash.concat(enc.Utf8.parse("BcryptSimple")));

    return hash;
}

/**/
function Bcrypt(data: WordArray, salt: WordArray, loops: number): WordArray {
    const hash = SHA512(data);
    let saltReverse = salt.clone();
    reverse(saltReverse);

    // eslint-disable-next-line no-param-reassign
    salt = SHA512(salt.concat(hash));
    const hashReverse = SHA512(reverse(data));
    reverse(hashReverse);

    saltReverse = SHA512(hashReverse.clone().concat(saltReverse));
    reverse(saltReverse);

    let newhash = lib.WordArray.create();
    newhash.concat(salt);
    for (let i = 0; i < hash.sigBytes; i += 1) {
        const b1 = getByteFromWords(hash.words, i);
        const b2 = getByteFromWords(hashReverse.words, i);
        // eslint-disable-next-line no-bitwise
        newhash.concat(lib.WordArray.create([b1 << 24], 1));
        // eslint-disable-next-line no-bitwise
        newhash.concat(lib.WordArray.create([b2 << 24], 1));
    }
    newhash.concat(saltReverse);
    newhash = SHA512(newhash);

    for (let i = 0; i < loops; i += 1) {
        if (i % 2 !== 0) {
            reverse(newhash);
        }
        newhash = SHA512(newhash);
    }

    const b1 = lib.WordArray.create();
    const b2 = lib.WordArray.create();
    const b3 = lib.WordArray.create();
    const b4 = lib.WordArray.create();
    for (let i = 0; i < 16; i += 1) {
        const byte1 = getByteFromWords(newhash.words, i * 4);
        const byte2 = getByteFromWords(newhash.words, i * 4 + 1);
        const byte3 = getByteFromWords(newhash.words, i * 4 + 2);
        const byte4 = getByteFromWords(newhash.words, i * 4 + 3);
        // eslint-disable-next-line no-bitwise
        b1.concat(lib.WordArray.create([byte1 << 24], 1));
        // eslint-disable-next-line no-bitwise
        b2.concat(lib.WordArray.create([byte2 << 24], 1));
        // eslint-disable-next-line no-bitwise
        b3.concat(lib.WordArray.create([byte3 << 24], 1));
        // eslint-disable-next-line no-bitwise
        b4.concat(lib.WordArray.create([byte4 << 24], 1));
    }
    const a1 = b4;
    const a2 = reverse(b2);
    const a3 = reverse(b3);
    const a4 = b1;
    newhash = lib.WordArray.create().concat(a1).concat(a2).concat(a3).concat(a4);

    newhash = SHA512(newhash.concat(enc.Utf8.parse("Bcrypt")));
    return newhash;
}


function reverse(data: WordArray): WordArray {
    const len = Math.floor(data.sigBytes / 2);
    for (let si = 0; si < len; si += 1) {
        const {words} = data;
        const ei = data.sigBytes - 1 - si;
        const start = getByteFromWords(words, si);
        const end = getByteFromWords(words, ei);
        setByteToWords(words, si, end);
        setByteToWords(words, ei, start);
    }
    return data;
}

function getByteFromWords(words: number[], indexByte: number): number {
    const i = Math.floor(indexByte / 4);
    const j = indexByte % 4;
    let temp: number = words[i];
    if (j === 0) {
        // eslint-disable-next-line no-bitwise
        temp >>>= 24;
    } else if (j === 1) {
        // eslint-disable-next-line no-bitwise
        temp = (temp & 0x00ff0000) >>> 16;
    } else if (j === 2) {
        // eslint-disable-next-line no-bitwise
        temp = (temp & 0x0000ff00) >>> 8;
    } else if (j === 3) {
        // eslint-disable-next-line no-bitwise
        temp &= 0x000000ff;
    }
    return temp;
}

function setByteToWords(words: number[], indexByte: number, num: number) {
    const i = Math.floor(indexByte / 4);
    const j = indexByte % 4;
    const twords = words;
    if (j === 0) {
        // eslint-disable-next-line no-bitwise
        twords[i] = (twords[i] & 0x00ffffff) | (num << 24);
    } else if (j === 1) {
        // eslint-disable-next-line no-bitwise
        twords[i] = (twords[i] & 0xff00ffff) | (num << 16);
    } else if (j === 2) {
        // eslint-disable-next-line no-bitwise
        twords[i] = (twords[i] & 0xffff00ff) | (num << 8);
    } else if (j === 3) {
        // eslint-disable-next-line no-bitwise
        twords[i] = (twords[i] & 0xffffff00) | num;
    }
}

const data = enc.Utf8.parse("hello world");
const key = enc.Utf8.parse("123");
console.log("BcryptSimple:", BcryptSimple(data, key, 2048).toString());
console.log("Bcrypt:", Bcrypt(data, key, 2048).toString());
