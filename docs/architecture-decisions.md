# Architecture Decision Records

## ADR-001: GitOps 도구로 ArgoCD 채택 (3장)
**시점**: 2026-08 / **결정**: GitOps 도구로 ArgoCD를 채택하고, Flux·Jenkins X·Spinnaker는 사용하지 않는다.
**이유**:
- Web UI로 배포 상태를 실시간 확인할 수 있어, 학습 과정에서 "지금 무슨 일이 일어나는지" 눈으로 볼 수 있다
- Application CRD로 "어떤 Git 경로 → 어떤 네임스페이스"를 선언적으로 관리한다
- selfHeal 기능으로 누군가 `kubectl edit`으로 직접 수정해도 Git 상태로 자동 복구된다
- GKE Standard의 e2-medium 노드에서 감당 가능한 리소스(~500MB 메모리)로 구동된다

## ADR-002: CI 도구로 GitHub Actions 채택 (3장)
**시점**: 2026-09 / **결정**: CI 도구로 GitHub Actions를 채택하고, Cloud Build·GitLab CI·Jenkins는 사용하지 않는다.
**이유**:
- 코드 저장소와 CI가 같은 GitHub 플랫폼에 있어 별도 서버 설치/관리가 불필요하다
- `.github/workflows/ci.yaml` 한 파일로 파이프라인을 선언적으로 정의한다
- 퍼블릭 저장소는 무제한, 프라이빗도 월 2,000분 무료 크레딧을 제공한다
- `google-github-actions/auth` 액션으로 GCP 인증 연동이 간편하다

## ADR-003: CI의 GCP 인증 방식으로 Workload Identity Federation 채택 (3장)
**시점**: 2026-09 / **결정**: CI→GCP 인증을 Service Account JSON 키 대신 Workload Identity Federation(WIF)으로 구현한다.
**이유**:
- 프로젝트 조직 정책(`constraints/iam.disableServiceAccountKeyCreation`)이 SA 키 발급 자체를 차단해, JSON 키 방식은 애초에 선택지가 아니었다
- 장기 유효한 키 파일을 GitHub Secrets에 저장할 필요가 없어 키 유출 위험이 구조적으로 사라진다
- GitHub Actions의 OIDC 토큰을 GCP IAM과 직접 연동해, 실행마다 단기 토큰만 발급받는다
- Workload Identity Pool/Provider/IAM 바인딩을 ch2에서 이미 준비해둬서 전환 비용이 낮았다

## ADR-004: 메트릭 모니터링으로 Prometheus + Grafana 채택 (4장)
**시점**: 2026-09 / **결정**: 메트릭 수집·시각화 도구로 Prometheus + Grafana(kube-prometheus-stack)를 채택하고, Datadog·CloudWatch·Google Cloud Monitoring은 사용하지 않는다.
**이유**:
- CNCF Graduated 프로젝트로, Kubernetes 모니터링의 사실상 표준이다
- SaaS 구독료 없이 자체 호스팅해 비용이 들지 않는다
- kube-prometheus-stack Helm 차트가 Prometheus·Grafana·Alertmanager 등 6개 컴포넌트를 검증된 버전 조합으로 한 번에 설치한다
- Grafana를 이후 Loki(로그)·Tempo(트레이스)와 하나로 통합해 도구가 파편화되지 않는다

## ADR-005: 외부 트래픽 관리로 Gateway API 채택 (5장)
**시점**: 2026-09 / **결정**: 외부 트래픽 라우팅 방식으로 Gateway API를 채택하고, Ingress NGINX·Istio·Traefik은 사용하지 않는다.
**이유**:
- Ingress를 대체하는 Kubernetes 공식 차세대 표준 API다(GA since K8s 1.27)
- GKE에서 네이티브로 지원해 별도 Controller 설치 없이 바로 사용 가능하다(`gke-l7-regional-external-managed`)
- Gateway(인프라팀 관리)와 HTTPRoute(앱팀 관리)로 역할이 분리되어 책임 경계가 명확하다
- HTTPRoute의 backendRefs로 트래픽을 분배할 수 있어, 5.3 Blue/Green 및 이후 Canary 전략과 연동된다

## ADR-006: 무중단 배포 전략으로 Argo Rollouts(Blue/Green) 채택 (5장)
**시점**: 2026-09 / **결정**: 무중단 배포 전략 도구로 Argo Rollouts를 채택하고 Blue/Green 전략을 적용하며, Flagger와 K8s 기본 Rolling Update는 사용하지 않는다.
**이유**:
- ArgoCD와 같은 Argo 생태계라 ArgoCD UI에서 Rollout 상태를 함께 확인할 수 있다
- CRD 기반 YAML 선언으로 배포 전략을 정의해 기존 GitOps 워크플로우와 호환된다
- Rollout CRD의 `strategy` 필드만 바꾸면 되므로, 5장 Blue/Green에서 6장 Canary로 점진적으로 전략을 진화시킬 수 있다
- `kubectl argo rollouts` 플러그인으로 preview→active 전환 등 배포 진행 상태를 실시간 모니터링할 수 있다
