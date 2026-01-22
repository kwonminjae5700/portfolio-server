# 이메일 인증 기능 구현 가이드

## 🎯 구현 완료 사항

### 백엔드 개발 완료 ✅

1. **모델 추가**: [`verification_code.go`](internal/models/verification_code.go)
   - 이메일 인증 코드를 저장하는 모델 생성
   - 만료 시간 및 사용 여부 추적

2. **이메일 유틸리티**: [`email.go`](internal/utils/email.go)
   - 6자리 랜덤 인증 코드 생성
   - SMTP를 통한 이메일 전송 기능

3. **설정 업데이트**: [`config.go`](internal/config/config.go)
   - SMTP 설정 추가 (Host, Port, From, Password)

4. **서비스 로직**: [`auth_service.go`](internal/services/auth_service.go)
   - `SendVerificationCode`: 인증 코드 생성 및 전송
   - `VerifyCode`: 인증 코드 검증

5. **API 엔드포인트**: [`auth_handler.go`](internal/handlers/auth_handler.go)
   - `POST /auth/send-verification-code`: 인증 코드 전송
   - `POST /auth/verify-code`: 인증 코드 검증

6. **라우트 연결**: [`routes.go`](internal/routes/routes.go)
   - 새로운 엔드포인트를 라우터에 등록

7. **데이터베이스 마이그레이션**: [`database.go`](internal/database/database.go)
   - VerificationCode 모델을 AutoMigrate에 추가

---

## 🚀 백엔드 설정 및 실행 방법

### 1. 환경 변수 설정

`.env` 파일을 생성하고 다음 내용을 추가하세요:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=portfolio_db

# Server Configuration
SERVER_PORT=8080
ENV=development

# JWT Configuration
JWT_SECRET=your-secret-key-change-this
JWT_EXPIRATION_HOURS=24

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=me@kwon5700.kr
SMTP_PASSWORD=your-gmail-app-password
```

### 2. Gmail 앱 비밀번호 설정

Gmail을 사용하는 경우, 앱 비밀번호를 생성해야 합니다:

1. Google 계정 관리 페이지로 이동
2. 보안 → 2단계 인증 활성화
3. 앱 비밀번호 생성
4. 생성된 16자리 비밀번호를 `SMTP_PASSWORD`에 입력

📖 상세 가이드: https://support.google.com/accounts/answer/185833?hl=ko

### 3. 데이터베이스 마이그레이션

```bash
# 마이그레이션 실행
make migrate
# 또는
go run cmd/migrate/main.go
```

### 4. 서버 실행

```bash
# 개발 모드 실행
make run
# 또는
go run cmd/server/main.go
```

---

## 📡 API 사용 예시

### 1. 인증 코드 전송

```bash
curl -X POST http://localhost:8080/auth/send-verification-code \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

**성공 응답**:

```json
{
  "message": "Verification code sent successfully"
}
```

**실패 응답 (이메일 중복)**:

```json
{
  "code": 409,
  "message": "Email already exists",
  "detail": "This email is already registered"
}
```

### 2. 인증 코드 검증

```bash
curl -X POST http://localhost:8080/auth/verify-code \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456"
  }'
```

**성공 응답**:

```json
{
  "message": "Email verified successfully"
}
```

**실패 응답 (잘못된 코드)**:

```json
{
  "code": 400,
  "message": "Invalid verification code",
  "detail": "Code not found or already used"
}
```

**실패 응답 (만료된 코드)**:

```json
{
  "code": 400,
  "message": "Verification code expired",
  "detail": "Please request a new code"
}
```

### 3. 회원가입 (인증 완료 후)

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "username",
    "password": "password123"
  }'
```

---

## 📋 프론트엔드 개발 가이드

프론트엔드 개발자는 [`docs/EMAIL_VERIFICATION_PRD.md`](docs/EMAIL_VERIFICATION_PRD.md) 문서를 참고하세요.

이 문서에는 다음 내용이 포함되어 있습니다:

- 사용자 플로우
- UI/UX 요구사항
- 화면 구성 및 컴포넌트
- API 엔드포인트 상세
- 상태 관리 가이드
- 에러 처리
- 테스트 케이스
- 구현 우선순위
- 권장 기술 스택

---

## 🔍 주요 기능 설명

### 인증 코드 생성 로직

- 암호학적으로 안전한 랜덤 6자리 숫자 생성
- 유효 시간: 10분
- 중복 방지: 이메일당 하나의 활성 코드만 유지

### 보안 기능

- 인증 코드는 한 번만 사용 가능
- 만료 시간 자동 체크
- 이메일 중복 가입 방지
- SMTP TLS/SSL 지원

### 이메일 템플릿

- HTML 형식의 깔끔한 이메일 디자인
- 6자리 코드를 강조하여 표시
- 유효 시간 안내
- 모바일 친화적인 반응형 디자인

---

## 🛠 트러블슈팅

### 이메일이 전송되지 않는 경우

1. **SMTP 설정 확인**
   - `SMTP_FROM`, `SMTP_PASSWORD`가 올바른지 확인
   - Gmail 앱 비밀번호를 사용하고 있는지 확인

2. **방화벽 설정**
   - SMTP 포트(587)가 열려있는지 확인
   - 회사 네트워크에서 SMTP가 차단되어 있지 않은지 확인

3. **Gmail 보안 설정**
   - 2단계 인증이 활성화되어 있는지 확인
   - "보안 수준이 낮은 앱 허용"이 아닌 "앱 비밀번호"를 사용

### 데이터베이스 오류

```bash
# 데이터베이스 초기화
make db-reset

# 마이그레이션 재실행
make migrate
```

---

## 📚 다음 단계

1. **프론트엔드 개발**
   - PRD 문서를 참고하여 UI 구현
   - API 연동
   - 상태 관리

2. **테스트**
   - 단위 테스트 작성
   - 통합 테스트
   - E2E 테스트

3. **개선 사항**
   - 이메일 전송 큐 시스템 (Redis + Celery)
   - 속도 제한 (Rate Limiting)
   - 이메일 템플릿 커스터마이징
   - 다국어 지원

---

## 💡 참고 자료

- [Golang SMTP 가이드](https://pkg.go.dev/net/smtp)
- [Gmail SMTP 설정](https://support.google.com/a/answer/176600?hl=ko)
- [GORM 문서](https://gorm.io/docs/)
- [Gin 프레임워크](https://gin-gonic.com/docs/)
