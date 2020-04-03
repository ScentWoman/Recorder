#!/bin/bash
#prerequisite: ffmpeg rclone lsof
dir="/home/Record/Save"
node="GD"

while :
do
    prefix=`date -d "6 hour ago" +"%Y-%m-%d"`
    for file in ${dir}/*.flv
    do
        if [[ -f "${file}" ]]; then
            isOpen=`lsof "${file}"`

            if [[ -z "${isOpen}" ]]; then
                if [[ -z $(/usr/bin/rclone ls ${node}:/Record/${prefix} 2>&1 | grep "directory not found") ]]; then
                    if [[ $(find "${file}" -type f -size +8M 2>/dev/null) ]]; then
                        /usr/bin/ffmpeg -i "${file}" -vcodec copy -acodec copy "${file%.flv}.mp4" && rm -rf "${file}"
                        /usr/bin/rclone copy -v --stats 10s "${file%.flv}.mp4" ${node}:/Record/${prefix}/ && rm -rf "${file%.flv}.mp4"
                    else
                        echo Delete trash file "${file}".
                        rm -rf "${file}"
                    fi
                else
                    /usr/bin/rclone mkdir ${node}:/Record/${prefix}
                    #/usr/bin/rclone mkdir ${node}:/Record/${prefix}/Backup
                fi
            fi
        fi
    done

    for file in ${dir}/Backup/*.flv
    do
        if [[ -f "${file}" ]]; then
            isOpen=`lsof "${file}"`

            if [[ -z "${isOpen}" ]]; then
                if [[ $(find "${file}" -type f -size +8M 2>/dev/null) ]]; then
                    /usr/bin/ffmpeg -i "${file}" -vcodec copy -acodec copy "${file%.flv}.mp4" && rm -rf "${file}"
                    /usr/bin/rclone copy -v --stats 10s "${file%.flv}.mp4" ${node}:/Record/${prefix}/Backup/ && rm -rf "${file%.flv}.mp4"
                else
                    echo Delete trash file "${file}".
                    rm -rf "${file}"
                fi
            fi
        fi
    done

    sleep 10
done
