package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) PullRequestCreate(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error) {
	// TODO: implement the business logic to create a pull request
	return models.PullRequest{}, nil
}

func (s *Service) PullRequestMerge(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error) {
	// TODO: implement the business logic to merge a pull request
	return models.PullRequest{}, nil
}

//   /pullRequest/reassign:
//     post:
//       tags: [PullRequests]
//       summary: Переназначить конкретного ревьювера на другого из его команды
//       requestBody:
//         required: true
//         content:
//           application/json:
//             schema:
//               type: object
//               required: [ pull_request_id, old_user_id ]
//               properties:
//                 pull_request_id: { type: string }
//                 old_user_id: { type: string }
//             example:
//               pull_request_id: pr-1001
//               old_reviewer_id: u2
//       responses:
//         '200':
//           description: Переназначение выполнено
//           content:
//             application/json:
//               schema:
//                 type: object
//                 required: [pr, replaced_by]
//                 properties:
//                   pr:
//                     $ref: '#/components/schemas/PullRequest'
//                   replaced_by:
//                     type: string
//                     description: user_id нового ревьювера
//               example:
//                 pr:
//                   pull_request_id: pr-1001
//                   pull_request_name: Add search
//                   author_id: u1
//                   status: OPEN
//                   assigned_reviewers: [u3, u5]
//                 replaced_by: u5
//         '404':
//           description: PR или пользователь не найден
//           content:
//             application/json:
//               schema: { $ref: '#/components/schemas/ErrorResponse' }
//         '409':
//           description: Нарушение доменных правил переназначения
//           content:
//             application/json:
//               schema: { $ref: '#/components/schemas/ErrorResponse' }
//               examples:
//                 merged:
//                   summary: Нельзя менять после MERGED
//                   value:
//                     error: { code: PR_MERGED, message: cannot reassign on merged PR }
//                 notAssigned:
//                   summary: Пользователь не был назначен ревьювером
//                   value:
//                     error: { code: NOT_ASSIGNED, message: reviewer is not assigned to this PR }
//                 noCandidate:
//                   summary: Нет доступных кандидатов
//                   value:
//                     error: { code: NO_CANDIDATE, message: no active replacement candidate in team }

func (s *Service) PullRequestReassign(ctx context.Context, prID string, oldUserID string) (models.PullRequest, string, error) {
	// TODO: implement the business logic to reassign a pull request reviewer
	return models.PullRequest{}, "", nil
}
