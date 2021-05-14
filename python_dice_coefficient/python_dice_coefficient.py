# Copyright (C) 2021 Alteryx, Inc. All rights reserved.
#
# Licensed under the ALTERYX SDK AND API LICENSE AGREEMENT;
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.alteryx.com/alteryx-sdk-and-api-license-agreement
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""Example pass through tool."""
from ayx_plugin_sdk.core import (
    InputConnectionBase,
    Plugin,
    ProviderBase,
    register_plugin,
    FieldType,
    RecordPacket,
)


def dice_coefficient(a, b):
    if a is None or b is None: return 0.0
    if not len(a) or not len(b): return 0.0
    """ quick case for true duplicates """
    if a == b: return 1.0
    """ if a != b, and a or b are single chars, then they can't possibly match """
    if len(a) == 1 or len(b) == 1: return 0.0

    """ use python list comprehension, preferred over list.append() """
    a_bigram_list = [a[i:i + 2] for i in range(len(a) - 1)]
    b_bigram_list = [b[i:i + 2] for i in range(len(b) - 1)]

    a_bigram_list.sort()
    b_bigram_list.sort()

    # assignments to save function calls
    lena = len(a_bigram_list)
    lenb = len(b_bigram_list)
    # initialize match counters
    matches = i = j = 0
    while (i < lena and j < lenb):
        if a_bigram_list[i] == b_bigram_list[j]:
            matches += 1
            i += 1
            j += 1
        elif a_bigram_list[i] < b_bigram_list[j]:
            i += 1
        else:
            j += 1

    score = float(2 * matches) / float(lena + lenb)
    return score


class python_dice_coefficient(Plugin):
    """A sample Plugin that passes data from an input connection to an output connection."""

    def __init__(self, provider: ProviderBase):
        """Construct the AyxRecordProcessor."""
        self.name = "Pass through"
        self.provider = provider
        self.text1 = provider.tool_config['Text1']
        self.text2 = provider.tool_config['Text2']
        self.outputField = provider.tool_config['OutputField']
        self.output_anchor = self.provider.get_output_anchor("Output")
        self.output_metadata = None

    def on_input_connection_opened(self, input_connection: InputConnectionBase) -> None:

        """Initialize the Input Connections of this plugin."""
        if input_connection.metadata is None:
            raise RuntimeError("Metadata must be set before setting containers.")

        input_connection.max_packet_size = 1000
        self.output_metadata = input_connection.metadata.clone()
        self.output_metadata.add_field(self.outputField, FieldType.double, source="Dice Coefficient (Python)")
        self.output_anchor.open(self.output_metadata)

    def on_record_packet(self, input_connection: InputConnectionBase) -> None:
        """Handle the record packet received through the input connection."""
        packet = input_connection.read()
        df = packet.to_dataframe()
        df[self.outputField] = df.apply(self.score_row, axis=1)
        out_packet = RecordPacket.from_dataframe(df=df, metadata=self.output_metadata)
        self.output_anchor.write(out_packet)

    def on_complete(self) -> None:
        """Handle for when the plugin is complete."""

    def score_row(self, row):
        return dice_coefficient(row[self.text1], row[self.text2])


AyxPlugin = register_plugin(python_dice_coefficient)
