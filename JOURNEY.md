# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ⬜ | | 이 저장소 세션에서는 미실행 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-08-30 | Artifact Registry Docker 인증(서울 리전) 포함 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-08-30 | notiflex-platform, public으로 생성 (기본값 private에서 변경) |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-08-30 | notiflex-cluster, Gateway API 활성화 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-08-30 | notiflex-api v0.1.0, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-08-31 | |
| ch3 | 3.2 GitOps 도구 | ⬜ | | |
| ch3 | 3.3 기능 추가 | ⬜ | | |
| ch3 | 3.4 CI | ⬜ | | |
| ch3 | 3.5 CI-CD 연결 | ⬜ | | |
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
| | | | |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-08-30 최초 설정 (ch6 valkey-go, ch8 OTel SDK 요구사항 대비) |
| Notiflex 이미지 | v0.1.0 | 2026-08-30 최초 빌드/배포 |
| ArgoCD | | |
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
