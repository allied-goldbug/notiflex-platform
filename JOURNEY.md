# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-09-01 | statusline 설정 포함 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-08-30 | Artifact Registry Docker 인증(서울 리전) 포함 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-08-30 | notiflex-platform, public으로 생성 (기본값 private에서 변경) |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-08-30 | notiflex-cluster, Gateway API 활성화 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-08-30 | notiflex-api v0.1.0, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-08-31 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-08-31 | ArgoCD 설치, notiflex-smb Application Synced/Healthy |
| ch3 | 3.3 기능 추가 | ✅ | 2026-08-31 | /version 엔드포인트 추가(v0.1.1) 후 문제로 롤백(v0.1.0) |
| ch3 | 3.4 CI | ✅ | 2026-09-01 | GitHub Actions, Workload Identity Federation 인증 (조직 정책상 SA 키 발급 차단) |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-09-01 | CI가 매니페스트 이미지 태그 자동 갱신 후 push, ArgoCD 자동 동기화까지 엔드투엔드 검증 완료 |
| ch4 | 4.2 메트릭 모니터링 | ⬜ | | |
| ch4 | 4.3 로그 수집 | ⬜ | | |
| ch4 | 4.4 알림 | ⬜ | | |
| ch5 | 5.2 트래픽 관리 | ⬜ | | |
| ch5 | 5.3 무중단 배포 | ⬜ | | |
| ch6 | 6.1 캐시 | ⬜ | | |
| ch6 | 6.2 시크릿 관리 | ⬜ | | |
| ch6 | 6.3 Canary 전환 | ⬜ | | |
| ch7 | 7.2 멀티 노드풀 | ⬜ | | |
| ch7 | 7.3 App of Apps | ⬜ | | |
| ch7 | 7.4 멀티테넌시 | ⬜ | | |
| ch8 | 8.1 메시징 | ⬜ | | |
| ch8 | 8.2 트레이싱 | ⬜ | | |
| ch8 | 8.3 CronJob | ⬜ | | |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps 도구 (3.2) | ArgoCD | (이전 세션 기록 없음, 클러스터 상태로 완료만 확인) | |
| CI GCP 인증 (3.4) | Workload Identity Federation | Service Account JSON 키 | 프로젝트 조직 정책(`constraints/iam.disableServiceAccountKeyCreation`)이 키 발급을 차단, 장기 키 유출 위험도 없는 WIF가 더 안전 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-08-30 최초 설정 (ch6 valkey-go, ch8 OTel SDK 요구사항 대비) |
| Notiflex 이미지 | sha-36b94bb | 2026-08-30 v0.1.0 최초 빌드/배포. v0.1.1(/version 엔드포인트) 시도 후 문제로 v0.1.0 롤백. 2026-09-01부터 CI가 git SHA 태그로 자동 관리 (3.5 CI-ArgoCD 연결 이후 시맨틱 버전 태그 대신 SHA 태그 사용) |
| ArgoCD | v3.5.2 | 2026-08-31 설치, notiflex-smb Application Synced/Healthy |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium (Spot) | 2 | notiflex-api (replicas: 2) |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| 2.5 | gcloud config의 프로젝트 ID(`minlife1217-gitaiops-project`)가 실제로는 존재하지 않아 클러스터 생성 불가 | 접근 가능한 `project-d64f9b5c-20c8-4906-95b`("My First Project")로 전환, Container API 활성화 후 진행 |
| 2.6 | `gcloud builds submit` 최초 실행 시 `PERMISSION_DENIED` — Cloud Build API를 막 활성화해 서비스 계정이 아직 없었고, 이 프로젝트는 빌더로 Compute Engine 기본 서비스 계정을 사용하도록 구성되어 있었음 | Compute 기본 서비스 계정에 `roles/storage.objectViewer`, `roles/artifactregistry.writer`, `roles/logging.logWriter` 부여 후 재시도하여 해결 |
| 3.4 | 로컬 kubectl 컨텍스트가 `gke-gcloud-auth-plugin` 미설치로 인증 실패, 저장소 CLAUDE.md에 적힌 kubectl 컨텍스트명(`gke-sysnet4admin_book_gitaiops`)도 실제 클러스터 컨텍스트명과 불일치 | `gcloud components install gke-gcloud-auth-plugin` 설치 및 PATH 등록, CLAUDE.md의 컨텍스트명을 실제 값(`gke_project-d64f9b5c-20c8-4906-95b_asia-northeast3-a_notiflex-cluster`)으로 수정 |
| 3.4 | CI SA의 JSON 키 발급 시도 시 조직 정책(`constraints/iam.disableServiceAccountKeyCreation`)으로 `FAILED_PRECONDITION` 발생 | Service Account 키 대신 Workload Identity Federation으로 인증 방식 전환 (Pool/Provider/IAM 바인딩은 ch2에서 이미 준비되어 있었음) |
| 3.4 | `gh` 저장소 push 시 `refusing to allow an OAuth App to create or update workflow` — 초기 `gh auth login` 토큰에 `workflow` 스코프가 없어 `.github/workflows/` 변경 push가 거부됨 | `gh auth refresh -h github.com -s workflow`로 스코프 추가 후 재시도 |
| 3.5 | 저장소 기본 workflow 권한이 `read`라서 CI가 매니페스트 커밋을 push할 수 없었음 (`GITHUB_TOKEN` 권한 부족) | `gh api -X PUT .../actions/permissions/workflow`로 `default_workflow_permissions=write` 설정 (독자 승인 후 진행) |
