/*
 *
 * Copyright 2023 puzzlewidgetserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package puzzlewidgetserver

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dvaumoron/puzzlegrpcserver"
	pb "github.com/dvaumoron/puzzlewidgetservice"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var errWidgetNotFound = errors.New("widget not found")
var errActionNotFound = errors.New("action not found")
var errInternal = errors.New("internal service error")

type Data = map[string]any
type ActionHandler = func(Data) (string, string, []byte, error)

type action struct {
	kind    pb.MethodKind
	path    string
	handler ActionHandler
}

type Widget map[string]action

// based on gin path convention, with the path "/view/:id/:name"
// the map passed to handler will contains "pathData/id" and "pathData/name" entries
// handler returned values are supposed to be redirect, templateName and data :
//  1. redirect is a redirect path (ignored if empty), to build an absolute one on the site the map contains the "CurrentUrl" entry
//  2. data could be :
//     - a json marshalled map which entries will be added to the data passed to the template engine with templateName
//     - or any raw data when the action kind is pb.MethodKind_RAW
func (w Widget) AddAction(actionName string, kind pb.MethodKind, path string, handler ActionHandler) {
	w[actionName] = action{kind: kind, path: path, handler: handler}
}

type widgetServerAdapter struct {
	pb.UnimplementedWidgetServer
	widgets map[string]Widget
	logger  *otelzap.Logger
}

func (s widgetServerAdapter) GetWidget(ctx context.Context, request *pb.WidgetRequest) (*pb.WidgetResponse, error) {
	widgetName := request.Name
	widget, ok := s.widgets[widgetName]
	if !ok {
		return nil, errWidgetNotFound
	}
	return &pb.WidgetResponse{Name: widgetName, Actions: convertActions(widget)}, nil
}

func (s widgetServerAdapter) Process(ctx context.Context, request *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	widget, ok := s.widgets[request.WidgetName]
	if !ok {
		return nil, errWidgetNotFound
	}
	action, ok := widget[request.ActionName]
	if !ok {
		return nil, errActionNotFound
	}

	var data Data
	if err := json.Unmarshal(request.Data, &data); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal json from call", zap.Error(err))
		return nil, errInternal
	}

	redirect, templateName, resData, err := action.handler(data)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to handle action", zap.Error(err))
		return nil, errInternal
	}
	return &pb.ProcessResponse{Redirect: redirect, TemplateName: templateName, Data: resData}, nil
}

type WidgetServer struct {
	inner   puzzlegrpcserver.GRPCServer
	widgets map[string]Widget
}

func Make(serviceName string, version string, opts ...grpc.ServerOption) WidgetServer {
	grpcServer := puzzlegrpcserver.Make(serviceName, version, opts...)
	return WidgetServer{inner: grpcServer, widgets: map[string]Widget{}}
}

func (s WidgetServer) CreateWidget(widgetName string) Widget {
	widget, ok := s.widgets[widgetName]
	if !ok {
		widget = Widget{}
		s.widgets[widgetName] = widget
	}
	return widget
}

func (s WidgetServer) Start() {
	pb.RegisterWidgetServer(s.inner, widgetServerAdapter{widgets: s.widgets, logger: s.inner.Logger})
	s.inner.Start()
}

func convertActions(widget Widget) []*pb.Action {
	actions := make([]*pb.Action, 0, len(widget))
	for key, value := range widget {
		actions = append(actions, &pb.Action{Kind: value.kind, Name: key, Path: value.path})
	}
	return actions
}
