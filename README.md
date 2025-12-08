# 🎓 Portfolio Board Server

> GoLang 기반 게시판 백엔드 서비스 - 커서 기반 무한 스크롤, 계층형 댓글, 재귀적 삭제 구현

## 📋 프로젝트 개요

이 프로젝트는 수행평가 요구사항에 따라 **기존 코드의 아쉬운 점을 개선**하고, **3가지 필수 심화 기술**을 모두 적용한 게시판 백엔드 서비스입니다.

### 🎯 개선 목표 및 적용 기술

#### **커서 기반 무한 스크롤 (Cursor-based Infinite Scroll)**

- **개선 목표**: 기존 OFFSET 방식의 페이지네이션은 데이터가 많아질수록 성능이 저하되는 문제가 있습니다.
- **적용 방법**: `last_id`를 활용한 커서 기반 페이지네이션으로 일관된 성능 보장
- **구현 위치**:
  - `internal/services/article_service.go` - `GetArticles()` 메서드
  - `internal/services/comment_service.go` - `GetCommentsByArticle()` 메서드
- **기술적 이점**:
  - O(1) 시간 복잡도로 일관된 조회 성능
  - 실시간 데이터 추가/삭제 시에도 중복/누락 없음
  - 모바일 환경의 무한 스크롤 UX에 최적화
  - 인덱스 활용으로 대용량 데이터에서도 빠른 조회

---

## 🏗️ 프로젝트 구조

```
portfolio-server/
├── cmd/
│   ├── server/          # 메인 서버 애플리케이션
│   └── migrate/         # 데이터베이스 마이그레이션 도구
├── internal/
│   ├── config/          # 설정 관리
│   ├── database/        # 데이터베이스 연결 및 마이그레이션
│   ├── models/          # 데이터 모델 (User, Article, Comment)
│   ├── services/        # 비즈니스 로직
│   ├── handlers/        # HTTP 핸들러
│   ├── middleware/      # JWT, CORS, 에러 핸들링
│   ├── errors/          # 커스텀 에러 정의
│   ├── utils/           # 유틸리티 함수
│   └── routes/          # 라우팅 설정
├── docker-compose.yml   # Docker 구성
├── Dockerfile          # 애플리케이션 이미지
├── Makefile           # 빌드 및 실행 스크립트
└── go.mod             # Go 모듈 의존성
```

---

## 🚀 주요 기능

### ✅ 필수 구현 기능

- [x] **회원가입/로그인** - JWT 기반 인증
- [x] **게시글 CRUD** - 생성, 조회, 수정, 삭제
- [x] **댓글 CRUD** - 댓글 시스템
- [x] **권한 검증** - 작성자만 수정/삭제 가능
- [x] **커서 기반 무한 스크롤** - 게시글 및 댓글 목록

### 🎨 심화 기능

- [x] **커서 기반 페이지네이션** - OFFSET 대신 ID 기반 커서 사용
- [x] **일관된 성능** - 데이터 증가에도 일정한 조회 속도
- [x] **조회수 증가** - 게시글 조회 시 자동 카운팅
- [x] **댓글 수 표시** - 게시글 목록에 댓글 수 포함

---

## 📡 API 명세

### 인증 (Authentication)

| Method | Endpoint         | Description | Auth Required |
| ------ | ---------------- | ----------- | ------------- |
| POST   | `/auth/register` | 회원가입    | ❌            |
| POST   | `/auth/login`    | 로그인      | ❌            |
| GET    | `/auth/profile`  | 프로필 조회 | ✅            |

### 게시글 (Articles)

| Method | Endpoint        | Description              | Auth Required |
| ------ | --------------- | ------------------------ | ------------- |
| GET    | `/articles`     | 게시글 목록 (무한스크롤) | ❌            |
| GET    | `/articles/:id` | 게시글 상세 조회         | ❌            |
| POST   | `/articles`     | 게시글 생성              | ✅            |
| PUT    | `/articles/:id` | 게시글 수정              | ✅            |
| DELETE | `/articles/:id` | 게시글 삭제              | ✅            |

### 댓글 (Comments)

| Method | Endpoint                         | Description            | Auth Required |
| ------ | -------------------------------- | ---------------------- | ------------- |
| GET    | `/articles/:article_id/comments` | 댓글 목록 (무한스크롤) | ❌            |
| POST   | `/comments`                      | 댓글 생성              | ✅            |
| PUT    | `/comments/:id`                  | 댓글 수정              | ✅            |
| DELETE | `/comments/:id`                  | 댓글 삭제              | ✅            |

### 🔍 커서 기반 무한 스크롤 사용 예시

```bash
# 첫 페이지 요청
GET /articles?limit=20

# 응답
{
  "articles": [...],
  "next_cursor": 45,  # 마지막 게시글의 ID
  "has_more": true
}

# 다음 페이지 요청
GET /articles?last_id=45&limit=20
```

---

## 🛠️ 로컬 실행 방법

### 사전 요구사항

- Go 1.21 이상
- Docker & Docker Compose (선택사항)

### 1. Docker Compose로 실행 (권장)

```bash
# .env 파일 생성
cp .env.example .env

# Docker Compose로 실행
docker-compose up -d

# 로그 확인
docker-compose logs -f app
```

서버가 `http://localhost:8080`에서 실행됩니다.

### 2. 로컬 환경에서 실행

```bash
# 의존성 설치
go mod download

# PostgreSQL 실행 (Docker)
docker-compose up -d postgres

# .env 파일 생성 및 수정
cp .env.example .env

# 데이터베이스 마이그레이션
make migrate
# 또는
go run cmd/migrate/main.go

# 서버 실행
make run
# 또는
go run cmd/server/main.go
```

### 3. Makefile 명령어

```bash
make help           # 사용 가능한 명령어 목록
make build          # 빌드
make run            # 실행
make test           # 테스트
make docker-up      # Docker 시작
make docker-down    # Docker 종료
make migrate        # 마이그레이션
```

---

## 🗄️ 데이터베이스 스키마

### Users

```sql
id, email, username, password, created_at, updated_at, deleted_at
```

### Articles

```sql
id, title, content, author_id, view_count, created_at, updated_at, deleted_at
```

### Comments

```sql
id, article_id, author_id, content, created_at, updated_at, deleted_at
```

---

## 🧪 테스트

```bash
# 모든 테스트 실행
make test

# 커버리지 포함
make test-coverage
```

---

## 📦 배포

### Docker 이미지 빌드

```bash
docker build -t portfolio-server .
```

### 환경 변수 설정

프로덕션 환경에서는 다음 환경 변수를 설정해야 합니다:

```env
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
JWT_SECRET=your-strong-secret-key
ENV=production
```

---

## 🔐 인증 방식

JWT (JSON Web Token) 기반 인증을 사용합니다.

**요청 헤더**:

```
Authorization: Bearer <your-jwt-token>
```

**토큰 유효기간**: 24시간 (설정 변경 가능)

---

## 📝 API 사용 예시

### 1. 회원가입

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### 2. 로그인

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. 게시글 작성

```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "제목",
    "content": "내용"
  }'
```

### 4. 게시글 목록 조회 (무한 스크롤)

```bash
# 첫 페이지
curl http://localhost:8080/articles?limit=20

# 다음 페이지
curl http://localhost:8080/articles?last_id=45&limit=20
```

### 5. 댓글 작성

```bash
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "article_id": 1,
    "content": "댓글 내용"
  }'
```

---

## 🎯 핵심 구현 사항

### 커서 기반 무한 스크롤

**게시글 목록**

- **위치**: `internal/services/article_service.go` - `GetArticles()` 메서드
- **핵심 로직**: `WHERE id < last_id ORDER BY id DESC LIMIT n+1`
- **장점**: 일관된 성능, 데이터 중복/누락 방지

**댓글 목록**

- **위치**: `internal/services/comment_service.go` - `GetCommentsByArticle()` 메서드
- **핵심 로직**: `WHERE article_id = ? AND id < last_id ORDER BY id DESC LIMIT n+1`
- **장점**: 대량의 댓글에서도 빠른 조회 성능

---

## 📚 기술 스택

- **언어**: Go 1.21
- **웹 프레임워크**: Gin
- **ORM**: GORM
- **데이터베이스**: PostgreSQL 15
- **인증**: JWT (golang-jwt/jwt)
- **컨테이너**: Docker & Docker Compose

---

## 👨‍💻 개발자

- **이름**: [권민재]
- **GitHub**: [@kwonminjae5700](https://github.com/kwonminjae5700)

---

## 📄 라이선스

MIT License

---

## 🙏 참고 자료

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [Cursor-based Pagination](https://www.sitepoint.com/paginating-real-time-data-cursor-based-pagination/)
