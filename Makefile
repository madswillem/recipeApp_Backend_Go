outfile:
	ffmpeg -i videos/in/big_buck_bunny_720p_surround.mp4 -c:v libx264 -preset veryfast -g 48 -keyint_min 48 -sc_threshold 0 -b:v 2500k -maxrate 2500k -bufsize 5000k -c:a aac -b:a 128k -hls_time 10 -hls_playlist_type vod -hls_segment_filename "videos/out/output%03d.ts" videos/out/playlist.m3u8

run_db:
	sudo systemctl start postgresql

run:
	go run main.go

test:
	go test -v -cover ./tests/*