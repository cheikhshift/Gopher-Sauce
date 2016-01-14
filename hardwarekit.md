 1. [Hardware Kit](#hardware-kit)
 2. [Requirements](#requirements)
 2. [Using your speakers](#using-your-speakers)
 3. [Taking a picture](#taking-a-picture)
 4. [Applinks](#applinks)
 5. [Notifications](#notifications)
 7. [Load Screen](#load-screen)
 8. [Launch Screen](#launch-screen)
 8. [Active JS](#active-js)
 6. [File management](#file-management)


# Hardware Kit
Hardware kit is an attempt to allow easy access to what we believe you need to make an incredible app. Hardware kit essentially is the bridge between your Go lang code and your device's capabilities.

Please keep in mind that in a Gopher Sauce application there is one instance of Hardware kit running; and that unique instance is available to methods that become API endpoints and templates.

# Requirements

- Gopher Sauce project with `bind` as its deploy target.
- Latest version of Python. 

### FYI 

.Layer  - Hardware kit instance within templates.
layer - Hardware kit instance within methods that become API endpoints. 

Here is the interface of the fabled Hardware kit : 

       type Flow interface {
         			PushView(url string)
         			DismissView()
         			DismissViewatInt(index int)
         			Width() float64
         			Height() float64
         			Device() int
         			ShowLoad()
         			HideLoad()
         			RunJS(line string)
         			Play(path string)
        			Stop()
        			SetVolume(power int)
        			GetVolume() int
        			IsPlaying() bool
        			PlayFromWebRoot(path string)
        			CreatePictureNamed(name string)
        			OpenAppLink(url string)
    
         			Notify(title string,message string)
         			AbsolutePath(file string) string
         			Download(url string, target string) bool
         			DownloadLarge(url string, target string)
         			Base64String(target string) string
         			GetBytes(target string) []byte
         			GetBytesFromUrl(target string) []byte
         			DeleteDirectory(path string) bool
         			DeleteFile(path string) bool
         			
         		}

# Using your speakers

The example below will play a sound from a template file. Templates will play song files within your web root. 

	{{ .Layer | PlaySound "/filename" }}

To stop the sound now :

	{{ .Layer | StopSound }}

To get or set the volume :

	{{ .Layer | SetVolume 40 }} // 40%
	{{ .Layer | GetVolume }}


Methods that become API endpoints have two ways to access your device's speaker.
		
	//Play from device's documents
	layer.Play("/filename")
	//load from compressed web root
	layer.PlayFromWebRoot("/filename")
	
	//layer set volume
	layer.SetVolume(40)
	power = layer.GetVolume()
	
# Taking a picture

To invoke your device's camera within your templates use the function below. One liner to take a picture and save it. 

	{{ .Layer | TakePicture "localNameOfImage" }}

Yup... it is that easy.

Within methods that become endpoints the call is as follows :

	layer.CreatePictureNamed("name")

# Applinks

App interconnectivity is the key to saving 5 days on reinventing a feature that an 'app-link' compatible app provides. The call is as follows within methods that become API endpoints... can't stress it anymore.

	layer.OpenAppLink("link...")

#Notifications

Deliver local notifications with one call, ideal for apps which utilizer background refresh to find new content. The call is as follows within templates :

	{{ .Layer | Notify "title" "Message" }}

Within methods that become endpoints please use this call : 

	layer.Notify("title", "Message")

# Load Screen

Don't ever leave your user thinking your app froze while you are attempting to handle a time consuming process. Load screen allows for adding a loading overlay spinner easy, while your application handles Big data. Here are the template methods :

	{{ .Layer | ShowLoad }}
	{{ .Layer | HideLoad }}

Within methods that become API endpoints :

		layer.ShowLoad()
		layer.HideLoad()
		
Cure the curiosity :)

![enter image description here](https://lh3.googleusercontent.com/-V7M5QZqcaXE/Vpf2t2yE_dI/AAAAAAAAAFU/f_WaAbJueH0/s0/image1.PNG "image1.PNG")

# Launch screen
To set a splash image for your application, simply put your image within your GoS project root and name it LaunchImage.png.

Yup... its that easy.

# Active JS

While your request is being processed it is possible to call the current active view's javascript with the Run method within templates. To run a javascript function within your templates please refer to the example below :

	{{ .Layer | Run "alert(\"nostalgia\")" }}

To run it within a method that is being used as an API endpoint :

	layer.RunJS("alert(\"\")")


# File management

This section covers accessing your app's own personal space within your future deployment target. Leverage your local filesystem to add more functionality to your app. 

### File management template methods

- Absolute Path : Find the full qualified path of a resource. Returns a string. ie: ` {{ .Layer | AbsolutePath "FILENAME" }}`
- Download : Downloads a file from the internet and saves it in the specified path. ie: ` {{ .Layer | Download "URL" "FILENAME" }}`
- Download Large : Made for files that are incredibly large, even opens a new thread. ie : ` {{.Layer | Download_lg "URL" "FILENAME" }}`
- Base64 : Will retrieve the base64 string of the specified file. ie : ` {{ .Layer | Base64 "FILENAME" }}`
- DeleteRes : Will delete the specified path, directory or file... ie : ` {{ .Layer | DeleteRes "FILENAME" }}`
