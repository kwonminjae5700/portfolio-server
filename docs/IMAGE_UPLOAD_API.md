# 이미지 업로드 API 명세서

## 개요

MinIO를 사용한 이미지 업로드 및 삭제 API

**Base URL**: `http://your-server-domain`  
**MinIO Endpoint**: `https://minio-api.kwon5700.kr`  
**Bucket**: `kwon5700-blog`

---

## 인증

모든 업로드 API는 JWT 인증이 필요합니다.

```
Authorization: Bearer {YOUR_JWT_TOKEN}
```

---

## API 엔드포인트

### 1. 이미지 업로드

**POST** `/upload/image`

이미지를 MinIO에 업로드하고 공개 URL을 반환합니다.

#### 요청

**Headers:**

```
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data
```

**Body (Form Data):**
| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| image | File | O | 업로드할 이미지 파일 |

**제한사항:**

- 최대 파일 크기: 10MB
- 허용 파일 형식: JPEG, JPG, PNG, GIF, WEBP

#### 응답

**Success (200 OK):**

```json
{
  "success": true,
  "data": {
    "url": "https://minio-api.kwon5700.kr/kwon5700-blog/images/550e8400-e29b-41d4-a716-446655440000-1706025600.jpg",
    "fileName": "550e8400-e29b-41d4-a716-446655440000-1706025600.jpg",
    "size": 245678
  }
}
```

**Error (400 Bad Request):**

```json
{
  "error": "이미지 파일을 찾을 수 없습니다"
}
```

```json
{
  "error": "파일 크기는 10MB를 초과할 수 없습니다"
}
```

```json
{
  "error": "허용되지 않는 파일 형식입니다. (jpeg, jpg, png, gif, webp만 가능)"
}
```

**Error (500 Internal Server Error):**

```json
{
  "error": "이미지 업로드에 실패했습니다"
}
```

#### 예제

**JavaScript (Fetch API):**

```javascript
const uploadImage = async (file, token) => {
  const formData = new FormData();
  formData.append("image", file);

  const response = await fetch("http://your-server/upload/image", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: formData,
  });

  return await response.json();
};

// 사용 예시
const fileInput = document.querySelector('input[type="file"]');
const file = fileInput.files[0];
const result = await uploadImage(file, "YOUR_JWT_TOKEN");
console.log(result.data.url); // 업로드된 이미지 URL
```

**JavaScript (Axios):**

```javascript
import axios from "axios";

const uploadImage = async (file, token) => {
  const formData = new FormData();
  formData.append("image", file);

  const response = await axios.post(
    "http://your-server/upload/image",
    formData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "multipart/form-data",
      },
    },
  );

  return response.data;
};
```

**cURL:**

```bash
curl -X POST http://your-server/upload/image \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "image=@/path/to/image.jpg"
```

---

### 2. 이미지 삭제

**DELETE** `/upload/image`

MinIO에서 이미지를 삭제합니다.

#### 요청

**Headers:**

```
Authorization: Bearer {JWT_TOKEN}
```

**Query Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| fileName | String | O | 삭제할 파일명 (UUID-timestamp.ext) |

#### 응답

**Success (200 OK):**

```json
{
  "success": true,
  "message": "이미지가 삭제되었습니다"
}
```

**Error (400 Bad Request):**

```json
{
  "error": "fileName이 필요합니다"
}
```

**Error (500 Internal Server Error):**

```json
{
  "error": "이미지 삭제에 실패했습니다"
}
```

#### 예제

**JavaScript (Fetch API):**

```javascript
const deleteImage = async (fileName, token) => {
  const response = await fetch(
    `http://your-server/upload/image?fileName=${fileName}`,
    {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  return await response.json();
};

// 사용 예시
await deleteImage(
  "550e8400-e29b-41d4-a716-446655440000-1706025600.jpg",
  "YOUR_JWT_TOKEN",
);
```

**cURL:**

```bash
curl -X DELETE "http://your-server/upload/image?fileName=550e8400-e29b-41d4-a716-446655440000-1706025600.jpg" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 파일명 규칙

업로드된 파일은 다음 형식으로 자동 생성됩니다:

```
{UUID}-{UNIX_TIMESTAMP}{EXTENSION}
```

예시: `550e8400-e29b-41d4-a716-446655440000-1706025600.jpg`

모든 이미지는 `images/` 폴더에 저장됩니다.

---

## MinIO 설정

**환경변수:**

```env
MINIO_ENDPOINT=minio-api.kwon5700.kr
MINIO_ACCESS_KEY=F122SHJRLSTH95AVUZYF
MINIO_SECRET_KEY=rMo8Q+ktptivV79PeWzDdqj10KJOQ1ZGNzVrsxZ8
MINIO_BUCKET=kwon5700-blog
MINIO_USE_SSL=true
```

---

## 에러 코드

| 상태 코드 | 설명                                             |
| --------- | ------------------------------------------------ |
| 200       | 성공                                             |
| 400       | 잘못된 요청 (파일 없음, 크기 초과, 형식 오류 등) |
| 401       | 인증 실패 (JWT 토큰 없음 또는 유효하지 않음)     |
| 500       | 서버 오류 (MinIO 업로드/삭제 실패)               |

---

## 주의사항

1. **파일 크기**: 최대 10MB까지만 업로드 가능
2. **파일 형식**: 이미지 파일만 허용 (JPEG, PNG, GIF, WEBP)
3. **인증 필수**: 모든 요청에 유효한 JWT 토큰 필요
4. **공개 URL**: 업로드된 이미지는 공개 URL로 접근 가능
5. **파일명 중복 방지**: UUID와 타임스탬프로 고유한 파일명 자동 생성
