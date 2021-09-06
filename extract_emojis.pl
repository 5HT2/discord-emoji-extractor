#!/usr/bin/perl

# Read from stdin
while (<STDIN>) {
    # Match the emoji regex, capturing each part of the emoji individually
    if (/<(a|):([A-z0-9_]+):([0-9]+)>/) {
        $suffix = ".png";
        # The a inside the emoji format <a:name:id> is only present if the emoji is a gif
        if ("$1" eq "a") {
            $suffix = ".gif";
        }
        
        # Save to the original file name and file extension
        # -nc will skip if the file exists (allowing resuming downloads)
        # The reason we do this instead of a --input-file flag is because of the file name
        print "wget -q --show-progress -nc -O $2$suffix https://cdn.discordapp.com/emojis/$3$suffix\n";
    }
}
