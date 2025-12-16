# 📚 Portfolio Server - 게시판 API 서버

> **Go + Gin + PostgreSQL + Docker**로 구축한 엔터프라이즈급 게시판 백엔드 서비스

---

## 🎯 프로젝트 개요

Portfolio Server는 단순한 CRUD 기능을 넘어 **실무 수준의 아키텍처 설계**를 적용한 게시판 API입니다.

### 핵심 기능

- ✅ **사용자 인증** (JWT 기반)
- ✅ **게시글 관리** (CRUD + 커서 기반 페이지네이션)
- ✅ **카테고리 관리** (Many-to-Many 관계)
- ✅ **댓글 시스템** (게시글별 댓글)
- ✅ **Swagger API 문서** (OpenAPI 기반)

### 기술 스택

- **언어**: Go 1.23
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL 15
- **ORM**: GORM
- **인증**: JWT (golang-jwt/jwt/v5)
- **컨테이너**: Docker & Docker Compose

---

## 🏗️ 아키텍처 설계

### 계층 분리 (Separation of Concerns)

```
┌─────────────────────────────────────────┐
│          HTTP Handler Layer             │  ← HTTP 요청/응답 처리
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│         Service Business Layer          │  ← 비즈니스 로직, 검증
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│      Database Access Layer (GORM)       │  ← 데이터 저장/조회
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│          Domain Models (GORM Tags)      │  ← 데이터 구조 정의
└─────────────────────────────────────────┘
```

### 디렉토리 구조

```
portfolio-server/
├── cmd/
│   └── server/main.go           # 애플리케이션 진입점
├── internal/
│   ├── config/config.go         # 환경설정 중앙 관리
│   ├── database/database.go     # DB 연결 & 마이그레이션
│   ├── errors/errors.go         # 구조화된 비즈니스 예외
│   ├── handlers/                # HTTP 핸들러 계층
│   ├── middleware/              # 미들웨어 (인증, 오류처리)
│   ├── models/                  # GORM 도메인 모델
│   ├── services/                # 비즈니스 로직 계층
│   ├── utils/password.go        # bcrypt 기반 비밀번호 관리
│   └── routes/routes.go         # 라우트 설정
├── docs/
│   ├── swagger.json             # OpenAPI 명세
│   └── swagger.html             # Swagger UI
├── docker-compose.yml           # 컨테이너 오케스트레이션
├── Dockerfile                   # 멀티 스테이지 빌드
├── go.mod                       # Go 모듈 정의
└── IMPROVEMENTS.md              # 개선사항 상세 문서
```

---

## 🔧 개선된 아키텍처 패턴

### 1️⃣ 중앙집중식 오류 처리 (Global Exception Handler)

**목적**: 모든 오류가 일관된 JSON 형식으로 응답되도록 보장

#### 구조화된 예외 정의

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 도메인별 예외 함수
func ErrUserNotFound() *AppError
func ErrArticleNotFound() *AppError
func ErrPermissionDenied() *AppError
```

#### 전역 에러 핸들러

```go
func ErrorHandler() gin.HandlerFunc {
    // 비즈니스 오류와 시스템 오류 자동 구분
    // Panic 발생 시 500 응답으로 자동 변환
}
```

**효과**:

- ✅ 일관된 HTTP 응답 형식
- ✅ 중앙 로깅으로 오류 추적 용이
- ✅ 서버 안정성 증가

---

### 2️⃣ 외부 설정 관리 (Externalized Configuration)

**목적**: 환경별로 다른 설정을 관리

#### 중앙집중식 설정

```go
func LoadConfig() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            // ... 환경변수로 로드
        },
    }
}
```

#### 환경변수 주입

```env
DB_HOST=localhost
DB_PORT=5432
JWT_SECRET=your-secret-key
ENV=development
```

**효과**:

- ✅ 환경별 설정 분리 (개발/테스트/운영)
- ✅ 보안 강화 (민감정보를 코드에서 제거)
- ✅ Docker 호환성 향상

---

### 3️⃣ 미들웨어 기반 코드 재사용 (Cross-Cutting Concerns)

**목적**: 인증, 오류 처리, CORS를 미들웨어로 일원화

#### 글로벌 미들웨어 적용

```go
router.Use(middleware.CORS())            // CORS 설정
router.Use(middleware.RecoveryHandler()) // Panic 처리
router.Use(middleware.ErrorHandler())    // 오류 처리

// 선택적 인증 미들웨어
articles.POST("", middleware.AuthMiddleware(), handler.CreateArticle)
```

**효과**:

- ✅ 코드 중복 제거
- ✅ 보안 정책 일원화
- ✅ 유지보수성 향상

---

### 4️⃣ 서비스 계층 분리 (Separation of Concerns)

**목적**: 비즈니스 로직을 서비스 계층에 집중

#### 핸들러는 HTTP만 담당

```go
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req services.CreateArticleRequest
    c.ShouldBindJSON(&req)              // HTTP 바인딩

    article, err := h.articleService.CreateArticle(...)  // 비즈니스 로직
    c.JSON(http.StatusCreated, article)  // HTTP 응답
}
```

#### 서비스는 비즈니스 로직 담당

```go
func (s *ArticleService) CreateArticle(req *CreateArticleRequest, authorID uint) (*Article, error) {
    // 검증, 도메인 로직, 카테고리 연결 등
    article := models.Article{...}
    s.db.Create(&article)

    // 복잡한 비즈니스 로직
    if len(req.CategoryIDs) > 0 {
        s.db.Model(&article).Association("Categories").Replace(categories)
    }
    return &article, nil
}
```

**효과**:

- ✅ 단위 테스트 용이
- ✅ 로직 재사용 가능
- ✅ 변경의 영향 범위 최소화

---

### 5️⃣ Many-to-Many 관계 설계 (Article-Category)

**목적**: 하나의 글이 여러 카테고리를 가질 수 있음

#### 카테고리 모델

```go
type Category struct {
    ID       uint      `gorm:"primaryKey"`
    Name     string    `gorm:"unique;not null"`
    Articles []Article `gorm:"many2many:article_categories;"`
}

// 조인 테이블
type ArticleCategory struct {
    ArticleID  uint `gorm:"primaryKey"`
    CategoryID uint `gorm:"primaryKey"`
}
```

#### 서비스에서 관계 관리

```go
// 글 작성 시 카테고리 자동 연결
if len(req.CategoryIDs) > 0 {
    var categories []models.Category
    s.db.Where("id IN ?", req.CategoryIDs).Find(&categories)
    s.db.Model(&article).Association("Categories").Replace(categories)
}
```

**효과**:

- ✅ 정규화된 데이터 구조
- ✅ 유연한 확장성
- ✅ 쿼리 효율성

---

## 📡 API 엔드포인트

### 인증 (Auth)

| Method | Endpoint         | 설명        | 인증 |
| ------ | ---------------- | ----------- | ---- |
| POST   | `/auth/register` | 회원가입    | ❌   |
| POST   | `/auth/login`    | 로그인      | ❌   |
| GET    | `/auth/profile`  | 프로필 조회 | ✅   |

### 게시글 (Articles)

| Method | Endpoint        | 설명                | 인증 |
| ------ | --------------- | ------------------- | ---- |
| GET    | `/articles`     | 글 목록 (커서 기반) | ❌   |
| GET    | `/articles/:id` | 글 상세 조회        | ❌   |
| POST   | `/articles`     | 글 작성             | ✅   |
| PUT    | `/articles/:id` | 글 수정             | ✅   |
| DELETE | `/articles/:id` | 글 삭제             | ✅   |

### 댓글 (Comments)

| Method | Endpoint                            | 설명      | 인증 |
| ------ | ----------------------------------- | --------- | ---- |
| GET    | `/articles/:id/comments`            | 댓글 목록 | ❌   |
| POST   | `/articles/:id/comments`            | 댓글 작성 | ✅   |
| PUT    | `/articles/:id/comments/:commentId` | 댓글 수정 | ✅   |
| DELETE | `/articles/:id/comments/:commentId` | 댓글 삭제 | ✅   |

### 카테고리 (Categories)

| Method | Endpoint          | 설명 | 인증 |
| ------ | ----------------- | ---- | ---- |
| GET    | `/categories`     | 목록 | ❌   |
| GET    | `/categories/:id` | 상세 | ❌   |
| POST   | `/categories`     | 생성 | ✅   |
| PUT    | `/categories/:id` | 수정 | ✅   |
| DELETE | `/categories/:id` | 삭제 | ✅   |

---

## 🚀 빠른 시작

### 전제 조건

- Docker & Docker Compose
- Go 1.23+ (로컬 개발 시)

### Docker로 실행

```bash
# 저장소 클론
git clone <repository-url>
cd portfolio-server

# Docker Compose로 실행
docker-compose up -d

# 서버 확인
curl http://localhost:8080/health

# Swagger 문서
open http://localhost:8080/swagger
```

### 로컬 개발 환경

```bash
# 의존성 설치
go mod download

# 환경 설정
cp .env.example .env

# PostgreSQL 실행 (Docker)
docker run --name postgres -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=portfolio_db -p 5432:5432 -d postgres:15-alpine

# 애플리케이션 실행
go run cmd/server/main.go
```

---

## 📝 API 사용 예시

### 1. 회원가입

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "username",
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

# 응답: {"token": "eyJhbGc...", "user": {...}}
```

### 3. 글 작성 (인증 필요)

```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{
    "title": "Go 언어 배우기",
    "content": "Go 언어는 동시성 처리에 강합니다...",
    "category_ids": [1, 2]
  }'
```

### 4. 글 목록 조회 (커서 기반 페이지네이션)

```bash
# 처음 조회
curl "http://localhost:8080/articles?limit=20"

# 다음 페이지
curl "http://localhost:8080/articles?limit=20&last_id=4"
```

### 5. 댓글 작성

```bash
curl -X POST http://localhost:8080/articles/5/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{"content": "좋은 글이네요!"}'
```

---

## 🔐 보안 기능

### JWT 기반 인증

- 토큰 만료 시간: 기본 24시간 (환경변수로 변경 가능)
- Bearer 토큰 방식: `Authorization: Bearer <token>`

### 비밀번호 보안

- bcrypt 알고리즘 (cost: 10)
- 평문 비밀번호는 저장 불가

### 권한 검증

- 글/댓글 수정/삭제는 작성자만 가능

---

## 📊 데이터베이스 스키마

### Many-to-Many 관계 (Article-Category)

```
┌─────────────┐         ┌──────────────────┐         ┌──────────────┐
│   Article   │◄───────►│ ArticleCategory  │◄───────►│  Category    │
│             │ 1    N  │  (조인 테이블)   │  N    1 │              │
└─────────────┘         └──────────────────┘         └──────────────┘
```

---

## 🐳 Docker 배포

### 멀티 스테이지 빌드

- 빌드 스테이지: golang:1.23-alpine
- 최종 이미지: alpine:latest
- 최종 크기: ~50MB

### 환경별 배포

```bash
# 개발 환경
docker-compose up

# 운영 환경
ENV=production docker-compose up -d
```

---

## 📈 성능 최적화

### 커서 기반 페이지네이션

- O(1) 시간복잡도
- 실시간 데이터 변화에 안전

### N+1 쿼리 방지

```go
s.db.Preload("Author").Preload("Categories").Find(&articles)
```

### 인덱싱

- author_id INDEX
- email UNIQUE INDEX
- username UNIQUE INDEX

---

## 🔄 개선 이력

- ✅ **중앙집중식 오류 처리** - Global Exception Handler 구현
- ✅ **외부 설정 관리** - 환경변수 기반 설정 분리
- ✅ **미들웨어 기반 재사용** - 인증, 오류 처리, CORS 통일
- ✅ **서비스 계층 분리** - 비즈니스 로직 독립화
- ✅ **Many-to-Many 관계** - 글-카테고리 유연한 관계 설계
- ✅ **API 문서화** - Swagger UI 제공
- ✅ **Docker 배포** - 멀티 스테이지 빌드로 최적화

**자세한 내용**: [IMPROVEMENTS.md](./IMPROVEMENTS.md) 참조

---

## 📚 추가 정보

- [프로젝트 개선사항 상세 분석](./IMPROVEMENTS.md)
- Swagger UI: http://localhost:8080/swagger
- API 명세: http://localhost:8080/docs/swagger.json

---

**마지막 업데이트**: 2025년 12월 16일
