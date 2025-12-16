# 포트폴리오 서버 - 개선사항 분석

## 🎯 프로젝트 개선 목표

이 프로젝트는 기본적인 게시판 API 구현에서 나아가, **실무 수준의 오류 처리**, **설정 관리**, **코드 재사용성**을 고려한 아키텍처를 적용하였습니다.

---

## 1️⃣ 개선사항 1: 중앙집중식 오류 처리 (Global Exception Handler)

### 🔴 기존 문제점

- 각 핸들러에서 개별적으로 오류를 처리하면 일관성 없는 HTTP 응답 형식 발생
- 비즈니스 로직 오류와 시스템 오류의 구분 불명확
- 오류 로그 관리 분산

### 🟢 개선 방안

#### 1.1 커스텀 비즈니스 예외 정의

📄 **파일**: `internal/errors/errors.go`

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 도메인별 예외 정의
func ErrUserNotFound() *AppError
func ErrArticleNotFound() *AppError
func ErrCommentNotFound() *AppError
func ErrPermissionDenied() *AppError
func ErrEmailAlreadyExists() *AppError
```

**특징**:

- ✅ 모든 오류가 구조화된 `AppError` 타입으로 통일
- ✅ HTTP 상태 코드가 예외 객체에 포함되어 일관성 보장
- ✅ 도메인별 오류 함수로 재사용성 높음

#### 1.2 전역 에러 핸들러 미들웨어

📄 **파일**: `internal/middleware/error.go`

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            // 비즈니스 오류와 시스템 오류 구분 처리
            if appErr, ok := err.(*errors.AppError); ok {
                c.JSON(appErr.Code, ErrorResponse{...})
                return
            }

            // 예상치 못한 오류는 500으로 통일
            c.JSON(http.StatusInternalServerError, ErrorResponse{...})
        }
    }
}

// Panic 복구도 중앙집중식으로 관리
func RecoveryHandler() gin.HandlerFunc { ... }
```

**적용 효과**:

- ✅ **일관된 응답 형식**: 모든 오류가 동일한 JSON 구조로 응답
- ✅ **중앙 로깅**: 예상치 못한 오류를 한 곳에서 로깅
- ✅ **Panic 방지**: 서버 크래시 방지 및 적절한 500 응답

#### 1.3 사용 예시

```go
// ❌ 나쁜 예: 개별 처리
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// ✅ 좋은 예: 전역 처리
if article == nil {
    return nil, errors.ErrArticleNotFound()
}
// 미들웨어가 자동으로 처리함
```

**HTTP 응답 통일**:

```json
{
  "code": 404,
  "message": "Article not found",
  "detail": "The requested article does not exist"
}
```

---

## 2️⃣ 개선사항 2: 외부 설정 관리 (Externalized Configuration)

### 🔴 기존 문제점

- 데이터베이스 URL, JWT 시크릿 등을 코드에 하드코딩하면 보안 문제
- 환경별(개발/테스트/운영) 설정 변경이 어려움
- 설정값 산재로 관리 복잡

### 🟢 개선 방안

#### 2.1 중앙집중식 설정 로더

📄 **파일**: `internal/config/config.go`

```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    JWT      JWTConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

func LoadConfig() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "portfolio_db"),
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
            ENV:  getEnv("ENV", "development"),
        },
        JWT: JWTConfig{
            Secret:          getEnv("JWT_SECRET", "your-secret-key"),
            ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
        },
    }
}
```

**구조화된 설정 계층**:

- ✅ **Database**: 데이터베이스 연결 정보 집중
- ✅ **Server**: 서버 포트, 환경 설정
- ✅ **JWT**: 인증 관련 설정

#### 2.2 환경 변수 기반 로드

📄 **파일**: `.env.example`

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=portfolio_db

# Server
SERVER_PORT=8080
ENV=development

# JWT
JWT_SECRET=your-secret-key-change-this
JWT_EXPIRATION_HOURS=24
```

**적용 효과**:

- ✅ **환경별 설정 분리**: 개발/테스트/운영 환경에서 `.env` 파일만 변경
- ✅ **보안 강화**: 민감한 정보(DB 비번, JWT 시크릿)를 코드에서 제거
- ✅ **Docker 호환성**: 컨테이너 환경에서 환경변수 주입 용이
- ✅ **기본값 제공**: 설정 미지정 시 합리적인 기본값 사용

#### 2.3 사용 예시

```go
// main.go
func main() {
    cfg := config.LoadConfig()

    // DB 연결
    database.InitDatabase(cfg)

    // JWT 초기화
    middleware.InitJWT(&cfg.JWT)

    // 서버 시작
    router.Run(":" + cfg.Server.Port)
}
```

---

## 3️⃣ 개선사항 3: 미들웨어 기반 코드 재사용성 (Cross-Cutting Concerns)

### 🔴 기존 문제점

- 모든 핸들러에서 인증, 오류 처리, 로깅을 반복 구현
- 보안 정책 변경 시 모든 핸들러 수정 필요
- 기능 추가/변경의 영향도 높음

### 🟢 개선 방안

#### 3.1 인증 미들웨어 재사용

📄 **파일**: `internal/middleware/jwt.go` + `internal/routes/routes.go`

```go
// 인증 미들웨어 정의 (한 곳)
func AuthMiddleware() gin.HandlerFunc { ... }

// 라우트에서 선택적으로 적용
articleHandler := handlers.NewArticleHandler()
articles := router.Group("/articles")
{
    articles.GET("", articleHandler.GetArticles)                     // 공개
    articles.POST("", AuthMiddleware(), articleHandler.CreateArticle) // 인증 필요
}
```

**재사용 효과**:

- ✅ **단일 책임 원칙**: 인증 로직은 미들웨어에만 존재
- ✅ **일관된 인증**: 모든 보호 라우트에서 동일한 검증 적용
- ✅ **보안 정책 변경 용이**: JWT → OAuth 변경 시 미들웨어만 수정

#### 3.2 오류 처리 미들웨어 재사용

```go
// 모든 라우트에 동일하게 적용
router.Use(middleware.ErrorHandler())
router.Use(middleware.RecoveryHandler())

// 핸들러는 오류만 반환하면 됨
// 미들웨어가 자동으로 HTTP 응답 생성
func (h *ArticleHandler) GetArticle(c *gin.Context) {
    article, err := h.articleService.GetArticleByID(id)
    if err != nil {
        c.Error(err)  // 미들웨어가 처리
        return
    }
    c.JSON(http.StatusOK, article)
}
```

**코드 재사용 효과**:

- ✅ **DRY 원칙**: 오류 처리 로직 중복 제거
- ✅ **핸들러 간소화**: 각 핸들러는 비즈니스 로직에만 집중
- ✅ **유지보수성**: 오류 처리 로직 변경 시 한 곳만 수정

---

## 4️⃣ 개선사항 4: 서비스 계층 분리 (Separation of Concerns)

### 아키텍처 계층 구분

```
Handler (HTTP 계층)
   ↓
Service (비즈니스 로직 계층)
   ↓
Database (데이터 접근 계층)
   ↓
Models (도메인 모델)
```

### 핸들러의 책임 축소

📄 **파일**: `internal/handlers/article_handler.go`

```go
// ✅ 핸들러는 HTTP 관련만 처리
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req services.CreateArticleRequest
    c.ShouldBindJSON(&req)           // HTTP 바인딩

    article, err := h.articleService.CreateArticle(...) // 비즈니스 로직은 위임

    c.JSON(http.StatusCreated, article) // HTTP 응답
}
```

### 서비스의 책임 강화

📄 **파일**: `internal/services/article_service.go`

```go
// ✅ 서비스는 비즈니스 로직 담당
func (s *ArticleService) CreateArticle(req *CreateArticleRequest, authorID uint) (*Article, error) {
    // 1. 데이터 검증
    // 2. 도메인 로직 실행
    // 3. 데이터 저장
    // 4. 카테고리 연결 등 복잡한 로직

    article := models.Article{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: authorID,
    }

    if err := s.db.Create(&article).Error; err != nil {
        return nil, fmt.Errorf("failed to create article: %w", err)
    }

    // 카테고리 연결 (비즈니스 로직)
    if len(req.CategoryIDs) > 0 {
        // ... 카테고리 처리
    }

    return &article, nil
}
```

**효과**:

- ✅ **단위 테스트 용이**: 서비스 계층을 독립적으로 테스트 가능
- ✅ **로직 재사용**: 서비스를 다른 인터페이스(CLI, gRPC)에서도 재사용
- ✅ **유지보수성**: 비즈니스 로직 변경이 HTTP 계층에 영향 없음

---

## 5️⃣ 개선사항 5: Many-to-Many 관계 설계 (Article-Category)

### 데이터 모델 개선

#### 5.1 조인 테이블 명시적 정의

📄 **파일**: `internal/models/category.go`

```go
// 기본 카테고리 모델
type Category struct {
    ID       uint      `gorm:"primaryKey"`
    Name     string    `gorm:"unique;not null"`
    Articles []Article `gorm:"many2many:article_categories;"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 조인 테이블 명시
type ArticleCategory struct {
    ArticleID  uint `gorm:"primaryKey"`
    CategoryID uint `gorm:"primaryKey"`
}
```

**설계 이점**:

- ✅ **정규화**: 중복 데이터 제거, 데이터 일관성 보장
- ✅ **유연성**: 하나의 글이 여러 카테고리 가능
- ✅ **쿼리 효율**: 필요한 데이터만 선택적으로 로드

#### 5.2 서비스 계층에서 관계 관리

```go
// 글 작성 시 카테고리 연결
func (s *ArticleService) CreateArticle(req *CreateArticleRequest, authorID uint) (*Article, error) {
    // ...

    // 카테고리 연결 (비즈니스 로직)
    if len(req.CategoryIDs) > 0 {
        var categories []models.Category
        if err := s.db.Where("id IN ?", req.CategoryIDs).Find(&categories).Error; err != nil {
            return nil, fmt.Errorf("failed to find categories: %w", err)
        }
        if err := s.db.Model(&article).Association("Categories").Replace(categories); err != nil {
            return nil, fmt.Errorf("failed to associate categories: %w", err)
        }
    }

    return &article, nil
}
```

---

## 📊 개선 효과 요약

| 개선사항             | 기존 문제             | 해결 방법                  | 효과                        |
| -------------------- | --------------------- | -------------------------- | --------------------------- |
| **전역 예외 처리**   | 일관성 없는 오류 응답 | ErrorHandler 미들웨어      | HTTP 응답 통일, 중앙 로깅   |
| **외부 설정 관리**   | 하드코딩된 설정 값    | config 패키지 + 환경변수   | 환경별 설정 분리, 보안 강화 |
| **미들웨어 재사용**  | 인증/로깅 반복 구현   | 미들웨어 기반 아키텍처     | 코드 중복 제거, 유지보수성  |
| **계층 분리**        | 혼재된 비즈니스 로직  | Handler/Service/Repository | 테스트 용이, 로직 재사용    |
| **M-to-M 관계 설계** | 카테고리 확장성 부족  | Many2Many + 조인 테이블    | 정규화, 유연한 데이터 구조  |

---

## 🧪 다음 개선 과제: 단위 테스트

현재 프로젝트에는 다음과 같은 개선이 필요합니다:

### 목표: 비즈니스 로직 100% 테스트 커버리지

#### 예시 1: 댓글 삭제 로직 테스트

```go
func TestDeleteCommentWithoutPermission(t *testing.T) {
    // Given: 다른 사용자의 댓글이 있음
    // When: 다른 사용자가 삭제 시도
    // Then: ErrPermissionDenied 오류 반환
}

func TestDeleteCommentSuccess(t *testing.T) {
    // Given: 작성자 본인이 댓글을 작성함
    // When: 삭제 요청
    // Then: 데이터베이스에서 제거됨
}
```

#### 예시 2: 글 작성 시 카테고리 연결 테스트

```go
func TestCreateArticleWithMultipleCategories(t *testing.T) {
    // Given: 여러 카테고리 ID 제공
    // When: 글 작성
    // Then: 모든 카테고리가 연결됨
}
```

---

## 결론

이 프로젝트는 **단순한 CRUD 기능**을 넘어:

- ✅ **엔터프라이즈급 오류 처리**
- ✅ **환경별 설정 관리**
- ✅ **코드 재사용성**
- ✅ **계층 분리 원칙**

등을 적용하여 **실무 수준의 아키텍처**를 구현했습니다.
