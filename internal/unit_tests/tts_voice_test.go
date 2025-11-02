package unittests

import (
	"context"
	"log"
	"testing"
	"time"

	"rio-go-model/configs"
	"rio-go-model/internal/services"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
)

var testStoryDB *database.StoryDatabase
var testDBInitialized bool

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("🔧 Initializing services at startup...")

	// Initialize global settings first
	configs.InitializeSettings()

	// Create and initialize AppService
	config := configs.LoadConfig()
	appService := services.NewAppService(config)

	// Initialize all services (Firestore, Storage, etc.)
	if err := appService.Initialize(ctx); err != nil {
		log.Printf("⚠️  AppService initialization failed (credentials may be missing): %v", err)
		log.Printf("ℹ️  Tests requiring database will be skipped")
		testStoryDB = nil
		testDBInitialized = false
		return
	}

	// Set the singleton instance
	services.SetInstance(appService)
	log.Println("✅ AppService singleton initialized successfully")

	// Get database instance from AppService
	testStoryDB = appService.GetFirestore()
	testDBInitialized = testStoryDB != nil
}

// TestGetVoiceList_Voices tests the GetVoiceList function
func TestGetVoiceList_Voices(t *testing.T) {
	tests := []struct {
		name     string
		voice    string
		wantNil  bool // true if we expect nil result
		minCount int  // minimum expected voices (if not nil)
	}{
		{
			name:     "Standard voice",
			voice:    "Standard",
			wantNil:  false,
			minCount: 1,
		},
		{
			name:     "Chirp3HD voice",
			voice:    "Chirp3HD",
			wantNil:  false,
			minCount: 1,
		},
		{
			name:     "Invalid voice",
			voice:    "InvalidVoice",
			wantNil:  true,
			minCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voiceList := util.GetVoiceList(tt.voice)

			if tt.wantNil {
				if voiceList != nil {
					t.Errorf("GetVoiceList(%q) = %v, want nil", tt.voice, voiceList)
				}
			} else {
				if voiceList == nil {
					t.Errorf("GetVoiceList(%q) = nil, want non-nil", tt.voice)
				} else if len(voiceList) < tt.minCount {
					t.Errorf("GetVoiceList(%q) returned %d voices, want at least %d", tt.voice, len(voiceList), tt.minCount)
				}
			}
		})
	}
}

// TestSuspendAudioAPI_Database tests the SuspendAudioAPI function
// Note: This requires a valid database connection and API trigger data
// This test will be skipped if database credentials are not available
func TestSuspendAudioAPI_Database(t *testing.T) {
	if !testDBInitialized || testStoryDB == nil {
		t.Skip("Skipping test: database not initialized (credentials may be missing)")
	}

	ctx := context.Background()

	// Test with "audio" API model
	suspended, tag, err := testStoryDB.SuspendAudioAPI(ctx, "audio")

	if err != nil {
		t.Errorf("SuspendAudioAPI() returned error: %v", err)
		return
	}

	// Verify return values are valid (not nil/empty if suspended)
	if suspended {
		if tag == "" {
			t.Errorf("SuspendAudioAPI() returned suspended=true but empty tag")
		}
		t.Logf("✅ API is suspended with tag: %s", tag)
	} else {
		t.Logf("✅ API is not suspended (tag: %s)", tag)
	}
}

// func TestLargeStoryAudioAPI(t *testing.T) {
// 	if !testDBInitialized || testStoryDB == nil {
// 		t.Skip("Skipping test: database not initialized (credentials may be missing)")
// 	}
// 	storyText := `ఒక ఉదయం, సూర్యుడు మిలమిల మెరుస్తుండగా, చిన్ని బుజ్జి అనే కుందేలు తన బొరియలో నుంచి బయటకు వచ్చింది. దాని చెవులు గాలికి రెపరెపలాడుతున్నాయి. "అమ్మో, ఎంత అందమైన రోజు!" అంది బుజ్జి, తన ముక్కును కదిలిస్తూ. దాని పక్కనే, నేల మీద, ఒక చిన్న బుడగ కనిపించింది. అది మెల్లగా పైకి లేస్తూ, రంగులు మారుతోంది. ఆకాశంలా నీలం, పువ్వులలా ఎరుపు, ఆకులులా పచ్చ. "ఓహ్, నువ్వెవరు?" అని అడిగింది బుజ్జి, ఆసక్తిగా. బుడగ మెల్లగా కిందకు దిగి, బుజ్జి ముందు ఆగింది. దాని లోపల నుంచి ఒక చిన్న, మెరిసే ఆకారం కనిపించింది. అది ఒక చిన్న మెదడులా ఉంది, కానీ దాని చుట్టూ ఒక వెలుగు ఉంది. "నేను ఆలోచనని," అంది ఆ ఆకారం, దాని స్వరం గంటల సవ్వడిలా ఉంది. "మీరు దేని గురించి ఆలోచిస్తున్నారు?" బుజ్జి ఆశ్చర్యంగా చూసింది. "ఆలోచన? అంటే ఏమిటి?" అని అడిగింది. మెదడులాంటి ఆకారం నవ్వింది. "ఆలోచన అంటే నీ తలలో జరిగే ఒక చిన్న మాయాజాలం. నువ్వు ఏదైనా చూసినప్పుడు, విన్నప్పుడు, అనుభూతి చెందినప్పుడు, నీ మెదడులో కొన్ని విద్యుత్ సంకేతాలు పరిగెడతాయి. అవి ఒకదానితో ఒకటి కలిసి కొత్త విషయాలను సృష్టిస్తాయి. అదే ఆలోచన." బుజ్జి కొంచెం తికమకపడింది. "విద్యుత్ సంకేతాలా? అవి ఎక్కడ నుంచి వస్తాయి?" అని ప్రశ్నించింది. "అవి నీ మెదడులోని చిన్న చిన్న కణాల నుంచి వస్తాయి. వాటిని న్యూరాన్స్ అంటారు. అవి ఒకదానితో ఒకటి మాట్లాడుకుంటాయి. నీ మెదడు ఒక పెద్ద కంప్యూటర్ లాంటిది, కానీ చాలా చాలా ఎక్కువ శక్తివంతమైనది." "వావ్!" అంది బుజ్జి, దాని కళ్ళు మెరిసాయి. "అంటే నేను నా మెదడుతో ఏదైనా చేయగలనా?" మెదడులాంటి ఆకారం ఉత్సాహంగా తలూపింది. "ఖచ్చితంగా! నువ్వు కొత్త విషయాలు నేర్చుకోవచ్చు, సమస్యలను పరిష్కరించవచ్చు, కలలు కనవచ్చు. నీ మెదడు ఒక అద్భుతమైన శక్తి." "కానీ కొన్నిసార్లు నాకు ఏమీ తోచదు," అంది బుజ్జి, కొంచెం బాధగా. "అప్పుడు నేను ఏం చేయాలి?" మెదడులాంటి ఆకారం బుజ్జి దగ్గరగా వచ్చి, దాని చెవిలో మెల్లగా చెప్పింది. "అప్పుడు కొత్త విషయాలు ప్రయత్నించు. ఒక కొత్త ఆట ఆడు, ఒక కొత్త పుస్తకం చదువు, స్నేహితులతో మాట్లాడు. కొత్త అనుభవాలు నీ మెదడుకు కొత్త ఇంధనం ఇస్తాయి. అవి కొత్త ఆలోచనలను పుట్టిస్తాయి." బుజ్జి ఆలోచనలో పడింది. "అంటే, నేను బయటకు వెళ్లి, ప్రపంచాన్ని చూస్తే, నా మెదడుకి ఎక్కువ శక్తి వస్తుందా?" మెదడులాంటి ఆకారం ఆనందంగా నవ్వింది. "అవును! నువ్వు ఎంత ఎక్కువ నేర్చుకుంటే, అంత ఎక్కువ ఆలోచించగలవు. నీ మెదడు ఒక కండరం లాంటిది. దాన్ని ఎంత వాడితే, అది అంత బలంగా మారుతుంది." అదే సమయం, ఒక ఉడుత చెట్టు మీద నుంచి దూకింది. "ఏం మాట్లాడుకుంటున్నారు మీరు?" అని అడిగింది. బుజ్జి ఉత్సాహంగా చెప్పింది, "నేను నా మెదడు గురించి నేర్చుకుంటున్నాను! అది చాలా శక్తివంతమైనది." ఉడుత ఆశ్చర్యంగా చూసింది. "నిజమా? అంటే నేను కూడా నా మెదడుతో గొప్ప పనులు చేయగలనా?" "తప్పకుండా!" అంది ఆలోచన. "మీరు ఒక చెట్టు కొమ్మ నుంచి మరొక కొమ్మకు దూకుతారు కదా? అది కూడా ఒక రకమైన ఆలోచనే. మీ మెదడు దూరం, వేగం, గాలిని అంచనా వేస్తుంది. అది ఒక గణితశాస్త్రజ్ఞుడిలా పనిచేస్తుంది!" బుజ్జి, ఉడుత ఒకరినొకరు చూసుకుని నవ్వుకున్నారు. అప్పటినుంచి, బుజ్జి తన మెదడును ఎప్పుడూ విస్మరించలేదు. అది కొత్త విషయాలను నేర్చుకోవడానికి, కొత్త పనులు చేయడానికి ఎప్పుడూ సిద్ధంగా ఉండేది. తన మెదడు ఒక అద్భుతమైన శక్తి అని దానికి తెలుసు. అది ఎప్పుడూ కొత్త ఆలోచనలతో నిండి ఉండేది.`
// 	if len(storyText) > 5000 {
// 		log.Println("Story text is too long, skipping test")
// 		return
// 	}

// 	tts := audio.NewGoogleTTS()
// 	fmt.Println("✅ Google TTS initialized successfully!")

// 	// Step 3: Generate audio
// 	fmt.Println("🔄 Step 3: Generating audio from SSML...")
// 	startTime := time.Now()

// 	audioData, tokens, err := tts.GenerateAudioAdapter(teluguText, "Telugu", "1", "Standard")
// 	if err != nil {
// 		log.Fatalf("❌ Failed to generate audio: %v", err)
// 	}

// 	storyDB.UpdateAPITokens(ctx, "audio", (int64)(tokens))
// 	if err != nil {
// 		log.Fatalf("❌ Failed to update api tokens: %v", err)
// 	}

// 	duration := time.Since(startTime)
// 	fmt.Printf("✅ Audio generated successfully in %v!\n", duration)
// 	fmt.Printf("📊 Audio data length: %d bytes\n", len(audioData))

// }
