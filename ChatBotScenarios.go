package IrisAPIs

import (
	"context"
	"fmt"
	ChatBot "github.com/rayer/chatbot"
	log "github.com/sirupsen/logrus"
)

type RootScenario struct {
	ChatBot.DefaultScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
	rs.RegisterState("entry", &EntryState{}, rs)
	rs.RegisterState("randomJoke", &RandomJokeState{}, rs)
	return nil
}

func (rs *RootScenario) EnterScenario(source ChatBot.Scenario) error {
	log.Debugln("Entering root scenario")
	return nil
}

func (rs *RootScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	log.Debugln("Exiting root scenario")
	return nil
}

func (rs *RootScenario) DisposeScenario() error {
	log.Debugln("Disposing root scenario")
	return nil
}

type EntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {
	es.Init(scenario, es)
	es.RegisterKeyword(&ChatBot.Keyword{Keyword: "system statistics", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
		err := es.InvokeNextScenario(&SystemStatisticScenario{}, ChatBot.Stack)
		return "Go to report scenario", err
	}})

	es.RegisterKeyword(&ChatBot.Keyword{Keyword: "random joke", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
		err := es.ChangeStateByName("randomJoke")
		return "Let me take a look for a joke....", err
	}})

}

func (es *EntryState) RawMessage() (string, error) {
	return "It is ChatBot demo! You can [check weather], browse [system statistics] or tell me a [random joke]", nil
}

type RandomJokeState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (ss *RandomJokeState) InitScenarioState(scenario ChatBot.Scenario) {
	ss.Init(scenario, ss)
	ss.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "next one",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			return "Here is another one...", nil
		},
	})
	ss.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			ss.ChangeStateByName("entry")
			return "Let's back to front door...", nil
		},
	})
}

func (ss *RandomJokeState) RawMessage() (string, error) {

	rj, err := fetchRandomJoke(context.TODO())
	var raw string
	if err != nil {
		raw = "Some bad happened : " + err.Error()
		return raw, err
	}
	raw = fmt.Sprintf("Here is a joke : %s\n%s\nInteresting in [next one]?\n", rj.Setup, rj.Punchline)
	//No keyword, no transform
	transformed, _, _ := ss.KeywordHandler.TransformRawMessage(raw)
	return transformed, nil
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}

// System Statistic
type SystemStatisticScenario struct {
	ChatBot.DefaultScenarioImpl
}

func (s *SystemStatisticScenario) InitScenario(uc *ChatBot.UserContext) error {
	s.RegisterState("Entry", &SystemStatisticEntryState{}, s)
	s.RegisterState("Docker", &DockerStatisticState{}, s)
	return nil
}

func (s *SystemStatisticScenario) EnterScenario(source ChatBot.Scenario) error {
	return nil
}

func (s *SystemStatisticScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	return nil
}

func (s *SystemStatisticScenario) DisposeScenario() error {
	return nil
}

func (s *SystemStatisticScenario) Name() string {
	return "SystemStatisticScenario"
}

type SystemStatisticEntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (s *SystemStatisticEntryState) InitScenarioState(scenario ChatBot.Scenario) {
	s.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.ExitScenario(scenario)
			return "Let's back...", nil
		},
	})
	s.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "overview",
		Action:  nil,
	})
}

func (s *SystemStatisticEntryState) RawMessage() (string, error) {
	return "Here are some sub module in system, you can [overview] them, or you can look in details of [system] or [docker] statistics", nil
}

type DockerStatisticState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (d *DockerStatisticState) InitScenarioState(scenario ChatBot.Scenario) {
	panic("implement me")
}

func (d *DockerStatisticState) RawMessage() (string, error) {
	panic("implement me")
}

type SystemOverviewState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (s *SystemOverviewState) InitScenarioState(scenario ChatBot.Scenario) {
	panic("implement me")
}

func (s *SystemOverviewState) RawMessage() (string, error) {
	panic("implement me")
}

func (s *SystemOverviewState) GetParentScenario() ChatBot.Scenario {
	panic("implement me")
}
