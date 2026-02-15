package obs

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Ccccraz/cogmoteGO/internal/commonTypes"
	"github.com/Ccccraz/cogmoteGO/internal/keyring"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/general"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/requests/stream"
	"github.com/andreykaipov/goobs/api/typedefs"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/spf13/viper"
)

// ObsStatus represents basic OBS server info.
type ObsStatus struct {
	ObsVersion string `json:"obs_version"`
	Streaming  bool   `json:"streaming"`
}

type InitObsResponse struct {
	SceneName     string `json:"scene_name"`
	SourceName    string `json:"source_name"`
	SceneFallback bool   `json:"scene_fallback"`
	SourceCreated bool   `json:"source_created"`
}

// obsData is the input payload for /obs/data
type obsData struct {
	MonkeyName  string  `json:"monkey_name"`
	TrialID     int64   `json:"trial_id"`
	StartTime   string  `json:"start_time"`
	CorrectRate float64 `json:"correct_rate"`
}

var client *goobs.Client
var obsProcess *os.Process

func strPtr(s string) *string { return &s }

func loadObsPassword() string {
	if password := strings.TrimSpace(os.Getenv("OBS_PASSWORD")); password != "" {
		return password
	}
	password, err := keyring.GetObsPassword()
	if err != nil {
		return ""
	}
	return password
}

// getDeviceName automatically retrieves hostname via gopsutil
func getDeviceName() string {
	hostInfo, err := host.Info()
	if err != nil {
		return "unknown"
	}
	return hostInfo.Hostname
}

func getObsCommand() (*exec.Cmd, error) {
	installMethod := viper.GetString("obs.install_method")
	switch installMethod {
	case "flatpak":
		return exec.Command("flatpak", "run", "com.obsproject.Studio"), nil
	case "system":
		switch runtime.GOOS {
		case "windows":
			return exec.Command("obs64.exe"), nil
		case "darwin":
			return exec.Command("open", "-a", "OBS"), nil
		default:
			return exec.Command("obs"), nil
		}
	default:
		return nil, fmt.Errorf("unknown obs install method: %s", installMethod)
	}
}

func startObsProcess() error {
	existingProc, err := findObsProcess()
	if err != nil {
		return fmt.Errorf("failed to check existing OBS process: %w", err)
	}
	if existingProc != nil {
		obsProcess = existingProc
		return nil
	}

	cmd, err := getObsCommand()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OBS: %w", err)
	}

	obsProcess = cmd.Process

	go func() {
		cmd.Wait()
		obsProcess = nil
	}()

	return nil
}

func findObsProcess() (*os.Process, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FI", "IMAGENAME eq obs64.exe", "/FO", "CSV", "/NH")
	case "darwin":
		cmd = exec.Command("pgrep", "-f", "OBS")
	default:
		cmd = exec.Command("pgrep", "-x", "obs")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return nil, nil
	}

	switch runtime.GOOS {
	case "windows":
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "obs64.exe") {
				fields := strings.Split(line, ",")
				if len(fields) >= 2 {
					pidStr := strings.Trim(strings.TrimSpace(fields[1]), `"`)
					pid := 0
					fmt.Sscanf(pidStr, "%d", &pid)
					if pid > 0 {
						return os.FindProcess(pid)
					}
				}
			}
		}
		return nil, nil
	default:
		pids := strings.Fields(strings.TrimSpace(string(output)))
		for _, pidStr := range pids {
			pid := 0
			fmt.Sscanf(pidStr, "%d", &pid)
			if pid > 0 && pid != os.Getpid() {
				return os.FindProcess(pid)
			}
		}
		return nil, nil
	}
}

func stopObsProcess() error {
	var proc *os.Process

	if obsProcess != nil {
		proc = obsProcess
	} else {
		var err error
		proc, err = findObsProcess()
		if err != nil {
			return fmt.Errorf("failed to find OBS process: %w", err)
		}
		if proc == nil {
			return nil
		}
	}

	if client != nil {
		client.Stream.StopStream(&stream.StopStreamParams{})
	}

	if runtime.GOOS == "windows" {
		if err := proc.Kill(); err != nil {
			return fmt.Errorf("failed to kill OBS: %w", err)
		}
		obsProcess = nil
		return nil
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to OBS: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		_, err := proc.Wait()
		done <- err
	}()

	select {
	case <-time.After(10 * time.Second):
		if err := proc.Kill(); err != nil {
			return fmt.Errorf("failed to kill OBS after timeout: %w", err)
		}
		obsProcess = nil
		return fmt.Errorf("OBS did not exit gracefully, force killed")
	case <-done:
		obsProcess = nil
		return nil
	}
}

func InitObsClient() (*InitObsResponse, error) {
	configSceneName := viper.GetString("obs.scene_name")
	sourceName := viper.GetString("obs.source_name")

	var opts []goobs.Option
	if password := loadObsPassword(); password != "" {
		opts = append(opts, goobs.WithPassword(password))
	}

	var err error
	client, err = goobs.New("localhost:4455", opts...)
	if err != nil {
		return nil, err
	}

	sceneListResp, err := client.Scenes.GetSceneList(&scenes.GetSceneListParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get scene list: %w", err)
	}

	if len(sceneListResp.Scenes) == 0 {
		return nil, fmt.Errorf("no scenes found in OBS")
	}

	var sceneName string
	var sceneFallback bool
	for _, s := range sceneListResp.Scenes {
		if s.SceneName == configSceneName {
			sceneName = configSceneName
			sceneFallback = false
			break
		}
	}
	if sceneName == "" {
		sceneName = sceneListResp.Scenes[0].SceneName
		sceneFallback = true
	}

	if sceneName == sourceName {
		return nil, fmt.Errorf("scene_name and source_name cannot be the same (both are '%s'): OBS does not allow a source to have the same name as a scene", sceneName)
	}

	sceneItems, err := client.SceneItems.GetSceneItemList(&sceneitems.GetSceneItemListParams{
		SceneName: strPtr(sceneName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get scene item list: %w", err)
	}

	for _, item := range sceneItems.SceneItems {
		if item.SourceName == sourceName {
			return &InitObsResponse{
				SceneName:     sceneName,
				SourceName:    sourceName,
				SceneFallback: sceneFallback,
				SourceCreated: false,
			}, nil
		}
	}

	data := obsData{
		MonkeyName:  "unknown",
		TrialID:     0,
		StartTime:   "unknown",
		CorrectRate: 0,
	}

	formatted := fmt.Sprintf(
		"%s %s %d %.2f%% %s",
		getDeviceName(),
		data.MonkeyName,
		data.TrialID,
		data.CorrectRate*100,
		data.StartTime,
	)

	var inputKind string
	switch runtime.GOOS {
	case "windows":
		inputKind = "text_gdiplus"
	case "linux", "darwin":
		inputKind = "text_ft2_source_v2"
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	var sceneItemId int

	createResp, err := client.Inputs.CreateInput(&inputs.CreateInputParams{
		SceneName: strPtr(sceneName),
		InputName: strPtr(sourceName),
		InputKind: strPtr(inputKind),
		InputSettings: map[string]any{
			"text": formatted,
		},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "ResourceAlreadyExists") {
			return nil, fmt.Errorf("failed to create input %s: %w", sourceName, err)
		}

		createItemResp, err := client.SceneItems.CreateSceneItem(&sceneitems.CreateSceneItemParams{
			SceneName:  strPtr(sceneName),
			SourceName: strPtr(sourceName),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add source %s to scene %s: %w", sourceName, sceneName, err)
		}
		sceneItemId = createItemResp.SceneItemId
	} else {
		sceneItemId = createResp.SceneItemId
	}

	_, err = client.SceneItems.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneName:   strPtr(sceneName),
		SceneItemId: &sceneItemId,
		SceneItemTransform: &typedefs.SceneItemTransform{
			PositionX:       0.0,
			PositionY:       1080.0,
			ScaleX:          1,
			ScaleY:          1,
			Alignment:       9,
			BoundsType:      "OBS_BOUNDS_SCALE_TO_HEIGHT",
			BoundsWidth:     1600.0,
			BoundsHeight:    60.0,
			BoundsAlignment: 9,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set scene item transform: %w", err)
	}

	return &InitObsResponse{
		SceneName:     sceneName,
		SourceName:    sourceName,
		SceneFallback: sceneFallback,
		SourceCreated: true,
	}, nil
}

func InitObsHandler(c *gin.Context) {
	resp, err := InitObsClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to initialize OBS client",
			Detail: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// HTTP handler: Get OBS status
func GetObsStatusHandler(c *gin.Context) {
	if client == nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "OBS client not initialized",
			Detail: "please call /obs/init first",
		})
		return
	}

	version, err := client.General.GetVersion(&general.GetVersionParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to get OBS version",
			Detail: err.Error(),
		})
		return
	}

	status, err := client.Stream.GetStreamStatus(&stream.GetStreamStatusParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to get OBS stream status",
			Detail: err.Error(),
		})
		return
	}

	resp := ObsStatus{
		ObsVersion: version.ObsVersion,
		Streaming:  status.OutputActive,
	}

	c.JSON(http.StatusOK, resp)
}

// HTTP handler: Start streaming
func PostStartObsStreamingHandler(c *gin.Context) {
	if client == nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "OBS client not initialized",
			Detail: "please call /obs/init first",
		})
		return
	}

	_, err := client.Stream.StartStream(&stream.StartStreamParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to start OBS stream",
			Detail: err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

// HTTP handler: Stop streaming
func PostStopObsStreamingHandler(c *gin.Context) {
	if client == nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "OBS client not initialized",
			Detail: "please call /obs/init first",
		})
		return
	}

	_, err := client.Stream.StopStream(&stream.StopStreamParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to stop OBS stream",
			Detail: err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func PostObsDataHandler(c *gin.Context) {
	if client == nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "OBS client not initialized",
			Detail: "please call /obs/init first",
		})
		return
	}

	var req obsData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commonTypes.APIError{
			Error:  "invalid request body",
			Detail: err.Error(),
		})
		return
	}

	formatted := fmt.Sprintf(
		"%s %s %d %.2f%% %s",
		getDeviceName(),
		req.MonkeyName,
		req.TrialID,
		req.CorrectRate*100,
		req.StartTime,
	)

	sourceName := viper.GetString("obs.source_name")
	_, err := client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: strPtr(sourceName),
		InputSettings: map[string]any{
			"text": formatted,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to set input settings",
			Detail: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func PostStartObsAppHandler(c *gin.Context) {
	if err := startObsProcess(); err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to start OBS",
			Detail: err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func PostStopObsAppHandler(c *gin.Context) {
	if err := stopObsProcess(); err != nil {
		c.JSON(http.StatusInternalServerError, commonTypes.APIError{
			Error:  "failed to stop OBS",
			Detail: err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func RegisterRoutes(r gin.IRouter) {
	obsGroup := r.Group("/obs")
	{
		obsGroup.GET("", GetObsStatusHandler)
		obsGroup.POST("/init", InitObsHandler)
		obsGroup.POST("/start", PostStartObsAppHandler)
		obsGroup.POST("/stop", PostStopObsAppHandler)
		obsGroup.POST("/data", PostObsDataHandler)
		obsGroup.POST("/streaming/start", PostStartObsStreamingHandler)
		obsGroup.POST("/streaming/stop", PostStopObsStreamingHandler)
	}
}
