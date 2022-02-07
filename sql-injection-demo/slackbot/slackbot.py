# Copyright 2018- The Pixie Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import time
import logging
import os
import schedule
import pxapi
from slack.web.client import WebClient
from slack.errors import SlackApiError


logging.basicConfig(level=logging.INFO)

# The slackbot requires the following configs, which are specified
# using environment variables. For directions on how to find these
# config values, see: https://docs.px.dev/tutorials/integrations/slackbot-alert/
if "PIXIE_API_KEY" not in os.environ:
    logging.error("Missing `PIXIE_API_KEY` environment variable.")
pixie_api_key = os.environ['PIXIE_API_KEY']

if "PIXIE_CLUSTER_ID" not in os.environ:
    logging.error("Missing `PIXIE_CLUSTER_ID` environment variable.")
pixie_cluster_id = os.environ['PIXIE_CLUSTER_ID']

if "SLACK_BOT_TOKEN" not in os.environ:
    logging.error("Missing `SLACK_BOT_TOKEN` environment variable.")
slack_bot_token = os.environ['SLACK_BOT_TOKEN']

# Slack channel for Slackbot to post in.
# Slack App must be a member of this channel.
if "SLACK_ALERT_CHANNEL" not in os.environ:
    logging.error("Missing `SLACK_ALERT_CHANNEL` environment variable.")
slack_channel = os.environ['SLACK_ALERT_CHANNEL']

pxl_script = open("sql_injections.pxl", "r").read()

def get_pixie_data(cluster_conn):
    msg = ["*Possible SQL injections detected in last 30 seconds*"]

    script = cluster_conn.prepare_script(pxl_script)

    # If you change the PxL script, you'll need to change the
    # columns this script looks for in the result table.
    for row in script.results("possible_sql_injections"):
        msg.append(format_message(row["rule_broken"],
                                  row["source"],
                                  row["req_body"]))

    if len(msg) == 1:
        return "*No SQL injections detected in last 30 seconds*"

    return "\n\n".join(msg)


def format_message(rule_broken, source, req_body):
    return (f"Rule `{rule_broken}` violated on `{source}`\n" +
        f"SQL: `{req_body}`")


def send_slack_message(slack_client, channel, cluster_conn):

    # Get data from the Pixie API.
    msg = get_pixie_data(cluster_conn)

    # Send a POST request through the Slack client.
    try:
        logging.info(f"Sending {msg!r} to {channel!r}")
        slack_client.chat_postMessage(channel=channel, text=msg)

    except SlackApiError as e:
        logging.error('Request to Slack API Failed: {}.'.format(e.response.status_code))
        logging.error(e.response)


def main():

    logging.debug("Authorizing Pixie client.")
    px_client = pxapi.Client(token=pixie_api_key)

    cluster_conn = px_client.connect_to_cluster(pixie_cluster_id)
    logging.debug("Pixie client connected to %s cluster.", cluster_conn.name())

    logging.debug("Authorizing Slack client.")
    slack_client = WebClient(slack_bot_token)

    # Send the first message right when we start up.
    send_slack_message(slack_client, slack_channel, cluster_conn)    

    # Schedule sending a Slack channel message every 30 seconds.
    schedule.every(30).seconds.do(lambda: send_slack_message(slack_client,
                                                            slack_channel,
                                                            cluster_conn))

    logging.info("Message scheduled for %s Slack channel.", slack_channel)

    while True:
        schedule.run_pending()
        time.sleep(5)


if __name__ == "__main__":
    main()
