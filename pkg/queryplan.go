package spanner

import (
	"context"
	"fmt"
	"io"
	"strings"

	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
)

// GetQueryPlan execute query plan
func (c *Client) GetQueryPlan(ctx context.Context, query string) (string, error) {
	// Create Session
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		c.config.Project, c.config.Instance, c.config.Database,
	)
	createSessionReq := &spannerpb.CreateSessionRequest{
		Database: databaseName,
	}
	session, err := c.spannerNomalClient.CreateSession(ctx, createSessionReq)
	if err != nil {
		return "", err
	}

	// Run Execute Sql Requeset
	req := &spannerpb.ExecuteSqlRequest{
		Session:   session.Name,
		Sql:       query,
		QueryMode: spannerpb.ExecuteSqlRequest_PLAN,
	}
	stream, err := c.spannerNomalClient.ExecuteStreamingSql(ctx, req)
	if err != nil {
		return "", err
	}
	var res string
	// Read response
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Get response node
		nodes := resp.GetStats().GetQueryPlan().PlanNodes

		// search node by index from Response Nodes
		searchPlanNode := func(index int32) *spannerpb.PlanNode {
			for _, node := range nodes {
				if node.Index == index {
					return node
				}
			}
			return nil
		}

		// make out put text by dfs
		indentlevel := 0
		var createOutput func(nodeId int32) string
		createOutput = func(nodeId int32) string {
			var result string
			node := searchPlanNode(nodeId)

			// 出力文を作成
			if node.DisplayName != "Reference" && node.DisplayName != "Function" && node.DisplayName != "Constant" && node.DisplayName != "Array Constructor" {
				result = strings.Repeat(" ", indentlevel)
				// Get Display Name
				if node.DisplayName == "Scan" {
					field := node.Metadata.Fields
					result += field["scan_type"].GetStringValue() + ":" + field["scan_target"].GetStringValue() + " "
				} else if node.DisplayName == "Distributed Union" {
					result += node.Metadata.Fields["call_type"].GetStringValue() + " " + node.DisplayName
				} else {
					result += node.DisplayName + ": "
				}

				// Get Content
				for _, nodeChild := range node.ChildLinks {
					searchednode := searchPlanNode(nodeChild.ChildIndex)
					if searchednode.DisplayName == "Reference" && node.DisplayName != "Serialize Result" {
						// Printout valiable matching
						result += nodeChild.Variable + ":=" + searchednode.ShortRepresentation.Description + " "
					} else if nodeChild.Type != "Split Range" {
						result += nodeChild.Type + " "
					}
				}

				result += "\n"
			}

			// recursion

			for _, nodeChild := range node.ChildLinks {
				indentlevel++
				result += createOutput(nodeChild.ChildIndex)
				indentlevel--
			}

			return result
		}

		res = createOutput(nodes[0].Index)
	}
	return res, nil
}
