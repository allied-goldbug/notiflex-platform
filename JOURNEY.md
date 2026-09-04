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
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-09-03 | kube-prometheus-stack 설치, Notiflex 대시보드(Pod CPU/메모리/재시작) 생성. HTTP 요청 수는 앱에 /metrics 미구현으로 보류 |
| ch4 | 4.3 로그 수집 | ✅ | 2026-09-04 | Loki(SingleBinary) + Fluent Bit(DaemonSet) 설치, Grafana에 Loki 데이터소스 추가. `{namespace="notiflex"}` 쿼리로 notiflex-api 로그 확인 완료 |
| ch4 | 4.4 알림 | ✅ | 2026-09-05 | PrometheusRule(PodRestartTooMany) + Alertmanager → Slack Webhook 연동. Webhook URL은 Secret `slack-webhook`(monitoring ns)으로만 저장, Git에는 커밋 안 함. 합성 알림으로 Slack 수신 검증 완료 |
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
| CI 도구 (3.4) | GitHub Actions | (이전 세션 기록 없음, 워크플로우 존재로 완료만 확인) | |
| 메트릭 모니터링 (4.2) | Prometheus + Grafana (kube-prometheus-stack) | Datadog, CloudWatch, Google Cloud Monitoring | 오픈소스 무료, K8s 사실상 표준, 이후 Loki/Tempo와 Grafana로 통합 |
| CI GCP 인증 (3.4) | Workload Identity Federation | Service Account JSON 키 | 프로젝트 조직 정책(`constraints/iam.disableServiceAccountKeyCreation`)이 키 발급을 차단, 장기 키 유출 위험도 없는 WIF가 더 안전 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-08-30 최초 설정 (ch6 valkey-go, ch8 OTel SDK 요구사항 대비) |
| Notiflex 이미지 | sha-36b94bb | 2026-08-30 v0.1.0 최초 빌드/배포. v0.1.1(/version 엔드포인트) 시도 후 문제로 v0.1.0 롤백. 2026-09-01부터 CI가 git SHA 태그로 자동 관리 (3.5 CI-ArgoCD 연결 이후 시맨틱 버전 태그 대신 SHA 태그 사용) |
| ArgoCD | v3.5.2 | 2026-08-31 설치, notiflex-smb Application Synced/Healthy |
| Loki | chart 7.3.0 (app 3.6.12) | 2026-09-04 설치, SingleBinary 모드, chunksCache/resultsCache 비활성화 |
| Fluent Bit | chart 2.6.0 (app v2.1.0) | 2026-09-04 설치, DaemonSet 2노드 모두 Running |
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
| 4.3 | `deploymentMode: SingleBinary`만 지정해도 `write`/`read`/`backend` 컴포넌트의 기본 replicas가 각 3이라 "more than zero replicas configured for both single binary and simple scalable targets" 에러로 설치 실패 | `helm-values/loki.yaml`에 `write.replicas: 0`, `read.replicas: 0`, `backend.replicas: 0` 명시 |
| 4.3 | Loki chart의 `chunksCache`/`resultsCache`(memcached 기반)가 기본 `enabled: true`이며 각각 500m CPU/9830Mi, 500m CPU/1229Mi를 요청 — 2노드 e2-medium에 `Insufficient cpu, Insufficient memory`로 Pending 발생 | `chunksCache.enabled: false`, `resultsCache.enabled: false`로 비활성화 (실습 규모에서는 캐시 불필요) |
| 4.3 | gz31 노드는 GKE 관리형 시스템 DaemonSet(kube-dns, fluentbit-gke, gke-metrics-agent 등)만으로 CPU requests가 이미 98%(940m 중 928m)에 도달해 있어, resource-budget.md의 Fluent Bit 권장치(25m/pod)로는 스케줄링 여유(12m)가 부족했음 | `helm-values/fluent-bit.yaml`의 CPU request를 5m으로 낮춰 두 노드 모두 Pending 없이 배치 성공. ch6 진입 전 resource-budget.md의 "위험 구간" 대응(Prometheus/Grafana/Alertmanager 축소)을 조금 더 앞당겨 검토할 필요 있음 |
| 4.3 | Loki 기본값 `auth_enabled: true`라 Fluent Bit/Grafana가 별도 `X-Scope-OrgID` 헤더 없이 요청하면 실패할 수 있음 | `loki.auth_enabled: false`로 단일 테넌트 모드 사용 (실습 규모에서는 멀티테넌시 불필요) |
| 4.3 | `commonConfig.replication_factor` 기본값이 3인데 `singleBinary.replicas: 1`과 맞지 않아 로그 쓰기 시 복제본 부족 문제가 될 수 있음 | `loki.commonConfig.replication_factor: 1`로 설정 |
| 4.3 | Grafana에 Loki 데이터소스가 자동으로 추가되지 않음 — 가드레일 문구의 "grafana.datasource.isDefault"는 loki chart에 존재하지 않는 키였음 | 실제로는 kube-prometheus-stack의 `grafana.additionalDataSources`(`helm-values/kube-prometheus.yaml`)에 Loki를 `isDefault: false`로 추가하고 `helm upgrade`로 반영 |
| 4.4 | 가드레일의 테스트 방법(`kubectl delete pod`)은 실제로 알림을 못 띄운다 — Pod을 지우면 ReplicaSet이 새 Pod을 만들 뿐이고 `kube_pod_container_status_restarts_total`은 새 Pod에서 0부터 시작하므로 "5분 내 재시작 2회 초과" 조건을 충족하지 못함 | Alertmanager API(`/api/v2/alerts`)로 `PodRestartTooMany` 합성 알림을 직접 전송해 Slack 수신 파이프라인만 검증. 실제 재시작 알림을 보려면 앱이 반복적으로 크래시하도록 만들어야 함 |
| 4.4 | Alertmanager 루트 receiver를 `null`→`slack-notifications`로 바꾸자, GKE가 컨트롤 플레인을 관리해 원래도 항상 거짓 양성으로 떠 있던 `TargetDown`/`KubeControllerManagerDown`/`KubeSchedulerDown`/`KubeProxyDown` 알림까지 한꺼번에 Slack으로 전송됨 (이전엔 `null` receiver라 조용히 무시되고 있었음) | `helm-values/kube-prometheus.yaml`에 `kubeControllerManager.enabled: false`, `kubeScheduler.enabled: false`, `kubeProxy.enabled: false`, `kubeEtcd.enabled: false` 추가 — GKE에서는 표준적으로 비활성화하는 항목들이라 해당 ServiceMonitor/알림 자체가 사라짐 |
